package services

import (
	"belimang/internal/utils"
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashingPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 7)
	return string(bytes), err
}

func ComparePassword(pass string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) == nil
}

type HashRes struct {
	hashing string
	err     error
}

type HashJob struct {
	usrPass string
	resChan chan HashRes
}

type HashingWorkerPool struct {
	jobs chan HashJob
}

func NewHashingPool(workerCount int, queueSize int) *HashingWorkerPool {
	h := &HashingWorkerPool{
	 jobs: make(chan HashJob, queueSize),
	}

	for i := 0; i < workerCount; i++ {
	 go h.worker()
	}

	return h
}

func (h *HashingWorkerPool) worker() {
	for job := range h.jobs {
		hashing, err := HashingPassword(job.usrPass)
		job.resChan <- HashRes{hashing: hashing, err: err}
		close(job.resChan)
	}
}

func (h *HashingWorkerPool) HashPasswordAsync(ctx context.Context, pass string) (string, error) {
	resultCh := make(chan HashRes, 1)
	job := HashJob{usrPass: pass, resChan: resultCh}

	select {
	case h.jobs <- job:
		select {
		case res := <-resultCh:
			return res.hashing, res.err
		case <-ctx.Done():
			return "", utils.NewTooManyReq("hashing canceled, server busy")
		case <-time.After(3 * time.Second):
			return "", utils.NewTooManyReq("hashing canceled, server busy")
		}
	default:
		return "", utils.NewTooManyReq("server busy, try again later")
	}
}
