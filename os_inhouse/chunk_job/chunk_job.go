// Package chunk_job provides a singleton buffered channel and a pool of
// workers (goroutines) to process chunkJob tasks concurrently.
// Use this to enqueue fixed-size “chunks” of data for background writing.
package chunk_job

import (
	"context"
	logger "docTrack/logger"
	"fmt"
	"os"
	"sync"
)

// -----------------------------------------------------------------------------
// chunkJob
// -----------------------------------------------------------------------------

// chunkJob represents a single unit of work: a chunk of bytes that needs
// to be written to disk (or elsewhere). All the information needed to
// write the chunk is carried in this struct.
type chunkJob struct {
	uploadID   string // Unique ID for the overall upload session
	parentPath string // Base directory path where chunk files should be stored
	chunkNO    uint   // Sequence number of this chunk within the upload
	data       []byte // Raw payload bytes for this chunk
}

func (job chunkJob) println() string {
	return fmt.Sprintf("chunk Job with upload ID %s and chunk no %d ", job.uploadID, job.chunkNO)
}

// -----------------------------------------------------------------------------
// Singleton buffered channel container
// -----------------------------------------------------------------------------

// bufferedChunkJobChannelStruct wraps our channel so we can attach methods
// and hide the raw channel behind a singleton API.
type bufferedChunkJobChannelStruct struct {
	jobs chan *chunkJob // buffered channel carrying chunkJob instances
}

var (
	// bufferedChunkJobChannelInstance holds the singleton instance.
	bufferedChunkJobChannelInstance *bufferedChunkJobChannelStruct

	// once ensures the singleton is only created once, even under concurrency.
	once sync.Once
)

// InstantiateBufferedChunkJobChannel initializes the singleton buffered channel
// with the specified capacity. Safe to call multiple times—only the first call
// actually creates the channel.
func InstantiateBufferedChunkJobChannel(bufferSize uint) {
	once.Do(func() {
		bufferedChunkJobChannelInstance = &bufferedChunkJobChannelStruct{
			jobs: make(chan *chunkJob, bufferSize),
		}
	})
}

// -----------------------------------------------------------------------------
// Worker pool management
// -----------------------------------------------------------------------------

// StartWorkerPool launches 'poolSize' goroutines that continuously read from
// the buffered channel and process each chunkJob. They listen on the provided
// context, exiting cleanly when ctx is cancelled.
//
// Returns an error if the channel hasn't been initialized yet.
func StartWorkerPool(ctx context.Context, poolSize uint) error {
	if bufferedChunkJobChannelInstance == nil {
		return fmt.Errorf("chunk job channel not initialized; call InstantiateBufferedChunkJobChannel first")
	}

	// Launch each worker goroutine.
	for i := uint(0); i < poolSize; i++ {
		go func(workerID uint) {
			// Worker loop: waits for either a new job or context cancellation.
			for {
				select {
				case <-ctx.Done():
					// Context was cancelled → exit this goroutine.
					return
				case job := <-bufferedChunkJobChannelInstance.jobs:
					// Received a job → handle the write.
					writeChunkAt(job)
				}
			}
		}(i)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Enqueue API
// -----------------------------------------------------------------------------

// AddChunkJob enqueues one chunkJob into the buffer. If the buffer is full,
// this call will block until a worker frees up space.
//
// Returns an error if the channel hasn't been initialized.
func AddChunkJob(job *chunkJob) error {
	if bufferedChunkJobChannelInstance == nil {
		return fmt.Errorf("chunk job channel not initialized; call InstantiateBufferedChunkJobChannel first")
	}
	bufferedChunkJobChannelInstance.jobs <- job
	return nil
}

// -----------------------------------------------------------------------------
// Consumer logic (stub)
// -----------------------------------------------------------------------------

// writeChunkAt is where you implement the actual I/O for a chunkJob.
// For example, you might:
//  1. Ensure job.parentPath exists (e.g. os.MkdirAll).
//  2. Construct a filename from job.uploadID and job.chunkNO.
//  3. Open/create the file.
//  4. Write job.data.
//  5. Handle errors or retries.
//
// Right now it’s a no-op stub that just prints the intended path.

// improvements I need to make
// pass a channel error channel
// if there is an error put it to that channel
func writeChunkAt(job *chunkJob, errChannel chan<- *chunkJob) {
	// get the filepath to the temporary chunk download
	filePath := fmt.Sprintf("%s/%s_chunk_%d.bin",
		job.parentPath, job.uploadID, job.chunkNO)

	// open a connection to the filePath
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		logger.ErrorLogger.Println(err)
		// need to add more error specific to the chunk job
		errChannel <- job
		return
	}
	// defer close the connection
	defer f.Close()
	// write to that connection
	_, err = f.Write(job.data)

	if err != nil {
		logger.ErrorLogger.Println(err)
		// need to add more error specific to the chunk job
		errChannel <- job
		return
	}

	logger.InfoLogger.Println(fmt.Sprintf("%s was succesfully writen ", job.println()))

}

// func writeChunkAt(path string, data []byte, _ int64) error {
// 	// Create or truncate the part file so it starts empty
// 	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()
// 	// Simply write the chunk; its file size will == len(data)
// 	_, err = f.Write(data)
// 	return err
// }
