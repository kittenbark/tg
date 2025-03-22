package tg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

// InputFile is either:
// - LocalFile for local files.
// - CloudFile for files in the cloud (i.e. by an id/url).
type InputFile interface {
	WriteMultipart(multipart *multipart.Writer, field string) error

	WriteMultipartAsAttachment(multipart *multipart.Writer) (value string, err error)
}

var (
	_ InputFile = (*LocalFile)(nil)
	_ InputFile = (*CloudFile)(nil)
)

// LocalFile This object represents the contents of a file to be uploaded.
// Must be posted using multipart/form-data in the usual way that files are uploaded via the browser.
type LocalFile struct {
	path string
}

func (file *LocalFile) WriteMultipartAsAttachment(multipart *multipart.Writer) (value string, err error) {
	field := strconv.FormatUint(rand.Uint64(), 16)
	if err := file.WriteMultipart(multipart, field); err != nil {
		return "", err
	}
	return "attach://" + field, nil
}

func (file *LocalFile) WriteMultipart(multipart *multipart.Writer, field string) error {
	fileWriter, err := multipart.CreateFormFile(field, filepath.Base(file.path))
	if err != nil {
		return err
	}

	fileReader, err := os.Open(file.path)
	if err != nil {
		return err
	}
	defer func() { _ = fileReader.Close() }()

	if _, err := io.Copy(fileWriter, fileReader); err != nil {
		return err
	}
	return nil
}

func FromDisk(path string) *LocalFile {
	return &LocalFile{path: path}
}

// CloudFile could be a telegram file_id or link to file (i.e. imgbb.com/abcde123).
type CloudFile string

func (i *CloudFile) WriteMultipartAsAttachment(multipart *multipart.Writer) (value string, err error) {
	return string(*i), nil
}

func (i *CloudFile) WriteMultipart(multipart *multipart.Writer, field string) error {
	return multipart.WriteField(field, string(*i))
}

func FromCloud(src string) *CloudFile {
	result := CloudFile(src)
	return &result
}

func GenericRequestMultipart[Request any, Result any](ctx context.Context, method string, request *Request) (result Result, err error) {
	token, err := tryGetTokenFromContext(ctx)
	if err != nil {
		return
	}
	url := fmt.Sprintf("%s/bot%s/%s", getOrDefault(ctx, ContextApiUrl, DefaultTelegramApiUrl), token, method)

	body, contentType := requestMultipartPreparePipes[Request](defaults(request))
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		_ = body.CloseWithError(err)
		return
	}
	httpRequest.Header.Set("Content-Type", contentType)

	httpResponse, err := getOrDefault(ctx, ContextHttpClient, http.DefaultClient).Do(httpRequest)
	if err != nil {
		return
	}
	defer func() { _ = httpResponse.Body.Close() }()

	type HttpResult struct {
		Ok          bool                   `json:"ok"`
		ErrorCode   int                    `json:"error_code,omitempty"`
		Description string                 `json:"description,omitempty"`
		Parameters  map[string]interface{} `json:"parameters,omitempty"`
		Result      Result                 `json:"result,omitempty"`
	}
	var httpResult HttpResult

	if err = json.NewDecoder(httpResponse.Body).Decode(&httpResult); err != nil {
		return
	}
	if !httpResult.Ok {
		err = newTelegramError(httpResult.ErrorCode, httpResult.Description, httpResult.Parameters)
		return
	}

	return httpResult.Result, nil
}

func multipartWritePipesInputMedia(media InputMedia, multipart *multipart.Writer) (string, error) {
	result := map[string]any{}

	requestValue := reflect.Indirect(reflect.ValueOf(media))
	requestType := requestValue.Type()
	for i := range requestValue.NumField() {
		fieldTag := requestType.Field(i).Tag.Get("json")
		if fieldTag == "-" {
			continue
		}
		fieldName, fieldJsonTagOptions, _ := strings.Cut(fieldTag, ",")

		fieldValue := requestValue.Field(i)
		if isEmptyValue(fieldValue) && strings.Contains(fieldJsonTagOptions, "omitempty") {
			continue
		}

		switch field := fieldValue.Interface().(type) {
		case InputFile:
			data, err := field.WriteMultipartAsAttachment(multipart)
			if err != nil {
				return "", err
			}
			result[fieldName] = data

		default:
			result[fieldName] = field
		}
	}

	data, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func multipartWritePipes[Request any](request *Request, pipe *io.PipeWriter, multipart *multipart.Writer) {
	defer func() {
		_ = multipart.Close()
		_ = pipe.Close()
	}()

	requestType := reflect.TypeFor[Request]()
	requestValue := reflect.Indirect(reflect.ValueOf(request))
	for i := range requestValue.NumField() {
		fieldTag := requestType.Field(i).Tag.Get("json")
		if fieldTag == "-" {
			continue
		}
		fieldName, fieldJsonTagOptions, _ := strings.Cut(fieldTag, ",")

		fieldValue := requestValue.Field(i)
		if isEmptyValue(fieldValue) && strings.Contains(fieldJsonTagOptions, "omitempty") {
			continue
		}

		switch field := fieldValue.Interface().(type) {
		case InputFile:
			if err := field.WriteMultipart(multipart, fieldName); err != nil {
				_ = pipe.CloseWithError(err)
				return
			}

		case InputMedia:
			media, err := multipartWritePipesInputMedia(field, multipart)
			if err != nil {
				_ = pipe.CloseWithError(err)
				return
			}
			if err := multipart.WriteField(fieldName, media); err != nil {
				_ = pipe.CloseWithError(err)
				return
			}

		case []InputMedia:
			medias := []string{}
			for _, inputMedia := range field {
				media, err := multipartWritePipesInputMedia(inputMedia, multipart)
				if err != nil {
					_ = pipe.CloseWithError(err)
					return
				}
				medias = append(medias, media)
			}
			if err := multipart.WriteField(fieldName, fmt.Sprintf("[%s]", strings.Join(medias, ","))); err != nil {
				_ = pipe.CloseWithError(err)
				return
			}

		case string:
			if err := multipart.WriteField(fieldName, field); err != nil {
				_ = pipe.CloseWithError(err)
				return
			}

		default:
			data, err := json.Marshal(field)
			if err != nil {
				_ = pipe.CloseWithError(err)
				return
			}
			if err := multipart.WriteField(fieldName, string(data)); err != nil {
				_ = pipe.CloseWithError(err)
				return
			}
		}
	}

	if err := multipart.Close(); err != nil {
		_ = pipe.CloseWithError(err)
		return
	}
}

func requestMultipartPreparePipes[Request any](request *Request) (reader *io.PipeReader, contentType string) {
	reader, writer := io.Pipe()
	multipartWriter := multipart.NewWriter(writer)

	// TODO: cancel io with context?
	go multipartWritePipes[Request](request, writer, multipartWriter)

	return reader, multipartWriter.FormDataContentType()
}

func GenericDownloadTemp(ctx context.Context, fileId string, dirAndPattern ...string) (filename string, err error) {
	tmp, err := os.CreateTemp(at(dirAndPattern, 0, ""), at(dirAndPattern, 1, ""))
	if err != nil {
		return "", err
	}
	defer func(name string) {
		if err != nil {
			_ = os.Remove(name)
		}
		if r := recover(); r != nil {
			_ = os.Remove(name)
			panic(r)
		}
	}(tmp.Name())

	if err := GenericDownload(ctx, tmp.Name(), fileId); err != nil {
		return "", err
	}
	return tmp.Name(), nil
}

func GenericDownload(ctx context.Context, path string, fileId string) error {
	file, err := GetFile(ctx, fileId)
	if err != nil {
		return err
	}

	return getOrDefault[downloaderFunc](ctx, ContextFileDownloadType, fileDownloadClassic)(ctx, file, path)
}

type downloaderFunc = func(ctx context.Context, file *File, path string) error

var (
	_ downloaderFunc = fileDownloadClassic
	_ downloaderFunc = fileDownloadLocalCopy
	_ downloaderFunc = fileDownloadLocalMove
)

func fileDownloadClassic(ctx context.Context, file *File, path string) error {
	token, err := tryGetTokenFromContext(ctx)
	if err != nil {
		return err
	}

	output, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(output *os.File) { _ = output.Close() }(output)

	url := fmt.Sprintf("%s/file/bot%s/%s", getOrDefault(ctx, ContextApiUrl, DefaultTelegramApiUrl), token, file.FilePath)
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := getOrDefault(ctx, ContextHttpClient, http.DefaultClient).Do(httpRequest)
	if err != nil {
		return err
	}
	defer func(body io.ReadCloser) { _ = body.Close() }(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram download file: bad status (%s)", resp.Status)
	}

	if _, err = io.Copy(output, resp.Body); err != nil {
		return err
	}

	return nil
}

func fileDownloadLocalCopy(ctx context.Context, file *File, path string) (err error) {
	source, err := os.Open(file.FilePath)
	if err != nil {
		return err
	}
	defer func() { _ = source.Close() }()

	destination, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = destination.Close()
		if err != nil {
			_ = os.Remove(destination.Name())
		}
	}()

	if _, err = io.Copy(destination, source); err != nil {
		return err
	}
	return nil
}

func fileDownloadLocalMove(ctx context.Context, file *File, path string) (err error) {
	defer func() {
		if err != nil {
			_ = os.Remove(file.FilePath)
			_ = os.Remove(path)
		}
	}()

	renameErr := os.Rename(file.FilePath, path)
	switch {
	case renameErr == nil:
		return nil

	case strings.Contains(renameErr.Error(), "invalid cross-device link"):
		if err := fileDownloadLocalCopy(ctx, file, path); err != nil {
			return err
		}
		return os.Remove(file.FilePath)

	default:
		return renameErr
	}
}
