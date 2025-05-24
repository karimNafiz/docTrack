// Package chunk_job provides a singleton buffered channel and a pool of
// workers (goroutines) to process chunkJob tasks concurrently.
// Use this to enqueue fixed-size “chunks” of data for background writing.
package chunk_job

import (
	"context"
	"fmt"
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

// -----------------------------------------------------------------------------
// Singleton buffered channel container
// -----------------------------------------------------------------------------

// bufferedChunkJobChannelStruct wraps our channel so we can attach methods
// and hide the raw channel behind a singleton API.
type bufferedChunkJobChannelStruct struct {
	jobs chan chunkJob // buffered channel carrying chunkJob instances
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
			jobs: make(chan chunkJob, bufferSize),
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
func AddChunkJob(job chunkJob) error {
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
func writeChunkAt(job chunkJob) {
	filePath := fmt.Sprintf("%s/%s_chunk_%d.bin",
		job.parentPath, job.uploadID, job.chunkNO)
	// TODO: replace with real filesystem write logic.
	fmt.Printf("Writing chunk to %q (size=%d bytes)\n", filePath, len(job.data))
}
