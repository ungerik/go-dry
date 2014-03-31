package dry

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"sync"
)

var (
	Deflate DeflatePool
	Gzip    GzipPool
)

// DeflatePool manages a pool of flate.Writer
// flate.NewWriter allocates a lot of memory, so if flate.Writer
// are needed frequently, it's more efficient to use a pool of them.
// The pool is unlimited and grows on demand.
// So if 100 writers are requested before
// any is returned, then 100 writers will be created.
type DeflatePool struct {
	writers []*flate.Writer
	mutex   sync.Mutex
}

// GetWriter returns flate.Writer from the pool, or creates a new one
// with flate.BestCompression if the pool is empty.
func (pool *DeflatePool) GetWriter(dst io.Writer) (writer *flate.Writer) {
	pool.mutex.Lock()

	count := len(pool.writers)
	if count == 0 {
		writer, _ = flate.NewWriter(dst, flate.BestCompression)
	} else {
		writer = pool.writers[count-1]
		pool.writers = pool.writers[:count-1]
		defer writer.Reset(dst) // do it after pool.mutex.Unlock()
	}

	pool.mutex.Unlock()

	return writer
}

// ReturnWriter returns a flate.Writer to the pool that can
// late be reused via GetWriter.
// Don't close the writer, Flush will be called before returning
// it to the pool.
func (pool *DeflatePool) ReturnWriter(writer *flate.Writer) {
	writer.Flush()
	pool.mutex.Lock()
	pool.writers = append(pool.writers, writer)
	pool.mutex.Unlock()
}

// Clean removes all flate.Writer from the pool.
func (pool *DeflatePool) Clean() {
	pool.mutex.Lock()
	pool.writers = pool.writers[0:0]
	pool.mutex.Unlock()
}

// GzipPool manages a pool of gzip.Writer.
// The pool is unlimited and grows on demand.
// So if 100 writers are requested before
// any is returned, then 100 writers will be created.
type GzipPool struct {
	writers []*gzip.Writer
	mutex   sync.Mutex
}

// GetWriter returns gzip.Writer from the pool, or creates a new one
// with gzip.BestCompression if the pool is empty.
func (pool *GzipPool) GetWriter(dst io.Writer) (writer *gzip.Writer) {
	pool.mutex.Lock()

	count := len(pool.writers)
	if count == 0 {
		writer, _ = gzip.NewWriterLevel(dst, gzip.BestCompression)
	} else {
		writer = pool.writers[count-1]
		pool.writers = pool.writers[:count-1]
		defer writer.Reset(dst) // do it after pool.mutex.Unlock()
	}

	pool.mutex.Unlock()

	return writer
}

// ReturnWriter returns a gzip.Writer to the pool that can
// late be reused via GetWriter.
// Don't close the writer, Flush will be called before returning
// it to the pool.
func (pool *GzipPool) ReturnWriter(writer *gzip.Writer) {
	writer.Flush()
	pool.mutex.Lock()
	pool.writers = append(pool.writers, writer)
	pool.mutex.Unlock()
}

// Clean removes all gzip.Writer from the pool.
func (pool *GzipPool) Clean() {
	pool.mutex.Lock()
	pool.writers = pool.writers[0:0]
	pool.mutex.Unlock()
}
