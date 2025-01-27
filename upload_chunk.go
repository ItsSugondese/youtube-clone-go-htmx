package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (u *DefaultUploader) UploadChunk(chunk ChunkMeta, totalChunk int, generatedFileName string) error {
	data, err := ioutil.ReadFile(chunk.FileName)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s?file_name=%s&chunk_name=%s&chunk_index=%d&total_chunks=%d", u.serverURL, generatedFileName, chunk.FileName, chunk.Index, totalChunk),
		bytes.NewReader(data),
	)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload chunk: %s", resp.Status)
	}

	return nil
}
