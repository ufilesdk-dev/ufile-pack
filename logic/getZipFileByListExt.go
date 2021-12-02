package logic

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"
	ufsdk "ufile-pack/gosdk"
	uflog "ufile-pack/gosdk/log"
	"ufile-pack/model"
)

type GetZipFileByListFileItem struct {
	Key    string `json:"key"`
	NewKey string `json:"new_key"`
}

type GetZipFileByListExtReq struct {
	Action     string                     `json:"action"`
	FileList   []GetZipFileByListFileItem `json:"file_list"`
	BucketName string                     `json:"bucket_name"`
	FileHost   string                     `json:"file_host"`
}

type GetZipFileByListExtRsp struct {
	Action   string                     `json:"Action"`
	FileList []GetZipFileByListFileItem `json:"file_list"`
	Key      string                     `json:"Key"`
	RetCode  int                        `json:"RetCode"`
	ErrMsg   string                     `json:"ErrMsg"`
}

func NewGetZipFileByListExtRsp() *GetZipFileByListExtRsp {
	return &GetZipFileByListExtRsp{
		Action:  "GetZipFileByListRsp",
		RetCode: 0,
		ErrMsg:  "",
	}
}

func GetZipFileByListExtRequest(msg []byte) ([]byte, error) {
	uflog.INFOF("GetZipFileByListExtRequest")
	var osr GetZipFileByListExtReq
	err := json.Unmarshal(msg, &osr)
	if err != nil {
		uflog.ERROR("GetUFileResourcePkg|json.Unmarshal|err:", err)
		return nil, err
	}

	uflog.DEBUG("request.body:", osr)
	osp := NewGetZipFileByListExtRsp()
	config := &ufsdk.Config{
		PublicKey:  model.G_Config.US3Config.PublicKey,
		PrivateKey: model.G_Config.US3Config.PrivateKey,
		BucketName: osr.BucketName,
		FileHost:   osr.FileHost,
	}
	req, err := getFileRequest(config)
	if err != nil {
		uflog.ERROR(err.Error())
		osp.ErrMsg = err.Error()
		osp.RetCode = -1
	} else {
		srcFile := osr.FileList
		destFile := "output/" + GenUuid() + ".zip"
		osp.Key = destFile
		osp.FileList = osr.FileList

		go PackFilesByListExt(req, srcFile, destFile)
	}
	response, err := json.Marshal(osp)
	if err != nil {
		uflog.ERROR("GetZipFileRequest|json.Marshal|err:", err)
		return nil, err
	}

	return response, nil
}

func PackFilesByListExt(req *ufsdk.UFileRequest, srcFileList []GetZipFileByListFileItem, destFile string) error {
	uflog.INFOF("create zip, srcFileList: %s, source_files: %s", srcFileList, destFile)

	var (
		partSise  = 1024 * 1024 * 4
		startTime = time.Now()
		queue     = make(chan KVExt)
		queuePart = make(chan *bytes.Buffer, 10)
		quit      = make(chan struct{})
	)
	// download
	download := func() {
		if len(srcFileList) != 0 {

			num := 0
			for _, content := range srcFileList {
				Key := strings.Replace(content.Key, " ", "", -1)
				NewKey := strings.Replace(content.NewKey, " ", "", -1)
				if Key == "" {
					continue
				}
				//var write bytes.Buffer
				rsp, err := req.DownloadFile(Key)
				if err != nil {
					uflog.ERRORF("DownloadFile err key = %s, err = %s ", Key, err.Error(), string(req.DumpResponse(true)))
					continue
				}
				data := KVExt{
					key:     Key,
					new_key: NewKey,
					data:    rsp,
				}
				queue <- data
				num++
			}

			close(queue)
		}
	}
	// compressed
	compressedToZip := func() {
		buf := new(bytes.Buffer)
		w := zip.NewWriter(buf)
		defer w.Close()
		num := 0
		for q := range queue {
			key, data, new_key := q.key, q.data, q.new_key
			uflog.INFOF("Download and Compression key = %s ", key)

			if new_key != "" { // key rename here
				key = new_key
			}

			f, err := w.Create(key)
			if err != nil {
				uflog.ERRORF("Create zipFile err %s", err.Error())
			}
			tmpBuf := make([]byte, 0, partSise*4)
			for {
				n, err := io.ReadFull(data.Body, tmpBuf[:cap(tmpBuf)])
				tmpBuf = tmpBuf[:n]
				if err != nil {
					if err == io.EOF {
						break
					}
					if err != io.ErrUnexpectedEOF {
						uflog.ERRORF("for-loop read data err", err)
						break
					}
				}
				if n == 0 {
					uflog.ERRORF("for-loop read data len = %d", n)
					break
				}
				//log.Println("read n bytes...", num, n)
				_, err = f.Write(tmpBuf)
				if err != nil {
					uflog.ERRORF("Write zipFile err %s", err.Error())
				}
				w.Flush()
				for buf.Len() >= partSise {
					tmp := make([]byte, partSise)
					buf.Read(tmp)
					queuePart <- bytes.NewBuffer(tmp)
				}
			}
			data.Body.Close()
			num++
			//uflog.INFOF("Compression complete key = %s", key)
		}
		//uflog.INFOF("Compression total nums = %d", num)
		err := w.Close()
		if err != nil {
			uflog.ERRORF("w.Close ", err)
		}
		if buf.Len() > 0 {
			uflog.DEBUG("last part len：", buf.Len())
			tmp := make([]byte, buf.Len())
			buf.Read(tmp)
			queuePart <- bytes.NewBuffer(tmp)
		}
		close(queuePart)
	}
	state, err := req.InitiateMultipartUpload(destFile, "application/x-zip-compressed")
	if err != nil {
		uflog.ERROR("InitiateMultipartUpload err", err, string(req.DumpResponse(true)))
		return err
	}
	if state != nil {
		partSise = state.BlkSize
	}
	go compressedToZip()
	go AsyncUpload(req, 10, state, destFile, queuePart, quit)

	download()
	<-quit
	endTime := time.Now()
	uflog.INFO("zipFiles cost time: ", endTime.Sub(startTime))
	return nil
}
