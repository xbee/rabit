package chunkstore

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type ChunkStore interface {
	Add(io.Reader, string) error
	LsFiles() ([]string, error)
	LsRemote() ([]string, error)
	CatFile(string, io.Writer) error
	Push(string) error
	Rm(string) error
	Fetch(string) error
	GC() error
	ChunkPath(string) string
}

type chunkStore struct {
	path   string
	remote string
}

func New(path, remote string) ChunkStore {
	return &chunkStore{path: path, remote: remote}
}

func (c *chunkStore) Add(r io.Reader, name string) error {
	w := newChunkWriter(c.path, r)
	spans, err := w.writeChunks(c)
	if err != nil {
		return err
	}

	return writeManifest(c.ManifestPath(name), spans)
}

func (c *chunkStore) CatFile(name string, w io.Writer) error {
	manifest, err := c.loadManifest(name)
	if err != nil {
		return err
	}

	for _, ch := range manifest.chunks {
		cpath := c.ChunkPath(ch)
		f, err := os.Open(cpath)
		if err != nil {
			return err
		}
		io.Copy(w, f)
		f.Close()
	}

	return nil
}

func (c *chunkStore) LsFiles() ([]string, error) {
	return nil, nil
}

func (c *chunkStore) LsRemote() ([]string, error) {
	return nil, nil
}

func (c *chunkStore) Push(name string) error {
	return nil
}

func (c *chunkStore) Rm(name string) error {
	return nil
}

func (c *chunkStore) Fetch(name string) error {
	return nil
}

func (c *chunkStore) GC() error {
	return nil
}

func (c *chunkStore) ChunkPath(hash string) string {
	prefix := hash[0:2]
	return filepath.Join(c.path, "chunks", prefix, hash)
}

func (c *chunkStore) ManifestPath(name string) string {
	return filepath.Join(c.path, "manifests", name)
}

func (c *chunkStore) loadManifest(name string) (*Manifest, error) {
	p := c.ManifestPath(name)
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	hashes := strings.Split(strings.TrimSpace(string(data)), "\n")
	return &Manifest{chunks: hashes}, nil
}

func writeManifest(path string, spans []span) error {
	var hashes []string
	for _, span := range spans {
		hashes = append(hashes, span.br)
	}

	manifest := Manifest{chunks: hashes}

	return ioutil.WriteFile(path, []byte(manifest.String()), 0660)
}
