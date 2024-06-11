package main

import "github.com/diegorezm/ds-fs-go/p2p"

type FileServerOpts struct {
	ListenAddr        string
	Transport         p2p.Transport
	StorageRoot       string
	PathTransformFunc PathTransformFunc
}

type FileServer struct {
	FileServerOpts
	store *Store
}

func NewFileServer(opts FileServerOpts) *FileServer {
	var storeOpts StoreOpts
	if len(opts.StorageRoot) > 0 {
		storeOpts.Root = opts.StorageRoot
	}
	if opts.PathTransformFunc != nil {
		storeOpts.PathTransformFun = opts.PathTransformFunc
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          &Store{storeOpts},
	}
}

func (s *FileServer) start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	return nil
}
