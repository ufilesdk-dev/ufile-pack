package logic

import (
	"bytes"
	"sync"
	ufsdk "ufile-pack/gosdk"
	uflog "ufile-pack/gosdk/log"
)

func AsyncUpload(u *ufsdk.UFileRequest, jobs int, state *ufsdk.MultipartState, keyName string, queue chan *bytes.Buffer, quit chan struct{}) error {
	//uflog.INFOF("AsyncUpload jobs = %d, BlkSize = %d, keyName = %s", jobs, state.BlkSize, keyName)
	var err error
	if jobs <= 0 {
		jobs = 1
	}
	if jobs >= 30 {
		jobs = 10
	}
	if state == nil {
		state, err = u.InitiateMultipartUpload(keyName, "application/x-zip-compressed")
		if err != nil {
			uflog.ERRORF("InitiateMultipartUpload err %s", err.Error())
			return err
		}
	}
	concurrentChan := make(chan error, jobs)
	for i := 0; i != jobs; i++ {
		concurrentChan <- nil
	}
	pos := 0
	wg := &sync.WaitGroup{}
	for b := range queue {
		//log.Printf("queue into data len = %d", len(b))
		uploadErr := <-concurrentChan //最初允许启动 10 个 goroutine，超出10个后，有分片返回才会开新的goroutine.
		if uploadErr != nil {
			err = uploadErr
			uflog.ERRORF("concurrentChan err %s", err.Error())
			break // 中间如果出现错误立即停止继续上传
		}
		wg.Add(1)
		go func(poss int, partData *bytes.Buffer) {
			defer wg.Done()
			dataLen := partData.Len()
			e := u.UploadPart(partData, state, poss)
			uflog.INFOF("Upload keyName: %s, PartId: %d, DataLen: %d", keyName, poss, dataLen)
			if e != nil {
				uflog.ERRORF("UploadPart err: %s", err)
			}
			concurrentChan <- e //跑完一个 goroutine 后，发信号表示可以开启新的 goroutine。
		}(pos, b)
		pos++
	}
	uflog.INFOF("wait for upload...... %d", pos)
	wg.Wait()       //等待所有任务返回
	if err == nil { //再次检查剩余上传完的分片是否有错误
	loopCheck:
		for {
			select {
			case e := <-concurrentChan:
				err = e
				if err != nil {
					break loopCheck
				}
			default:
				break loopCheck
			}
		}
	}
	close(concurrentChan)
	if err != nil {
		u.AbortMultipartUpload(state)
		return err
	}
	defer func() {
		uflog.INFOF("AsyncUpload file keyName = %s, partNum = %d", keyName, pos)
		quit <- struct{}{}
	}()
	err = u.FinishMultipartUpload(state)
	if err != nil {
		uflog.ERRORF("err:", err, string(u.DumpResponse(true)))
	}
	return nil
}

func GenUuid() string {
	return ufsdk.NewUUIDV4().String()
}
