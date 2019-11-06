package library

import (
	"github.com/golang/glog"
	"http/static/upload"
	"io"
	"myproject/dapi/o/fileupload"
	"myproject/dapi/x/math"
	"net/http"
	"os"
	f "path/filepath"
	"strconv"
	"strings"
)

const maxFileNameLen = 128

func makeSafeFilename(name string) string {
	v := upload.Slugify(name, true)
	if len(v) > maxFileNameLen {
		return v[:maxFileNameLen]
	}
	return v
}

func (s *FileUploadServer) HandleUploadFile(w http.ResponseWriter, r *http.Request) {

	//accept only POST method
	if r.Method != "POST" {
		http.Error(w, "METHOD NOT ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	file := fileupload.FileUpload{}

	fileName := r.URL.Query().Get("file_name")

	//file
	fileFromRequest, header, err := r.FormFile("file")

	//thumnail
	thumnailFromRequest, _, _ := r.FormFile("thumnail")

	if err != nil {
		w.Write([]byte("fail to read file form " + err.Error()))
		return
	}

	if len(file.Name) == 0 {
		fileName = header.Filename
	}

	file.Type = GetExtension(fileName)

	file.Name = fileName[0 : len(fileName)-len(f.Ext(fileName))]

	safeFileName := makeSafeFilename(fileName)

	file.PhysicalPath = f.Join(s.Dir, file.Type, safeFileName)

	file.Size = header.Size

	err = s.CreateDirIfNotExist(file.PhysicalPath)

	if err != nil {
		w.Write([]byte("fail to create filepath " + file.PhysicalPath + " error: " + err.Error()))
	}

	if _, err := os.Stat(file.PhysicalPath); !os.IsNotExist(err) {
		safeFileName = s.RandNewFileName(safeFileName)
		file.PhysicalPath = f.Join(s.Dir, file.Type, safeFileName)
	}

	file.RelativePath = strings.Replace(f.Join("/upload", file.Type, safeFileName), "\\", "/", -1)

	//write file to disk
	outstreamFile, err := os.Create(file.PhysicalPath)
	if err != nil {
		glog.Error("create ", file.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer outstreamFile.Close()

	if s.Limit == 0 {
		_, err = io.Copy(outstreamFile, fileFromRequest)
	} else {
		instream := &io.LimitedReader{N: s.Limit, R: fileFromRequest}
		_, err = io.Copy(outstreamFile, instream)
	}

	if err != nil {
		glog.Error("save", file.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if thumnailFromRequest != nil {
		//gen thumnail file name
		// thumnailFileName := file.Name
		file.PhysicalPathThumbnail = f.Join(s.Dir, file.Type, safeFileName+"_thumnail.png")
		file.RelativePathThumbnail = strings.Replace(f.Join("/upload/", file.Type, safeFileName+"_thumnail.png"), "\\", "/", -1)

		//write thumnail to disk

		outstreamThumnail, err := os.Create(file.PhysicalPathThumbnail)
		if err != nil {
			glog.Error("create ", file.Name, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer outstreamThumnail.Close()

		if s.Limit == 0 {
			_, err = io.Copy(outstreamThumnail, thumnailFromRequest)
		} else {
			instream := &io.LimitedReader{N: s.Limit, R: thumnailFromRequest}
			_, err = io.Copy(outstreamThumnail, instream)
		}
	}

	if err != nil {
		glog.Error("save", file.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = file.Create()

	if err != nil {
		glog.Error("save", file.Name, err)
		w.Write([]byte("create_file_to_db_fail"))
		return
	}

	// p := strings.Replace(f.Join("/upload", file.Type, file.Name), "\\", "/", -1)

	s.SendData(w, file)
}

func (s *FileUploadServer) HandleMarkDeleteFile(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("file_id")

	file, err := fileupload.GetByID(id)

	if err != nil {
		s.SendError(w, err)
		return
	}

	if file.DTime != 0 {
		s.Success(w)
		return
	}

	os.Remove(file.PhysicalPath)
	os.Remove(file.PhysicalPathThumbnail)

	err = fileupload.MarkDelete(id)

	if err != nil {
		s.SendError(w, err)
		return
	}

	s.Success(w)
}

func (s *FileUploadServer) RandNewFileName(filepath string) string {
	ext := f.Ext(filepath)
	filepath = strings.Trim(filepath, ext)
	filepath = filepath + "_" + math.RandString("", 5)
	return filepath + ext
}

func (s *FileUploadServer) CreateDirIfNotExist(path string) error {
	dir, _ := f.Split(path)
	return os.MkdirAll(dir, os.ModePerm)
}

func GetExtension(path string) string {
	ext := f.Ext(path)

	if len(ext) == 0 {
		return "other"
	}

	ext = strings.ToLower(strings.Trim(ext, "."))

	if ext == "jpg" || ext == "png" || ext == "gif" || ext == "jpeg" {
		return "image"
	}

	if ext == "mp4" {
		return "video"
	}

	if ext == "mp3" || ext == "flac" {
		return "audio"
	}

	return strings.Trim(ext, ".")
}

func (s *FileUploadServer) HandleGetByChapter(w http.ResponseWriter, r *http.Request) {
	pageSize := StrToInt(r.URL.Query().Get("page_size"))
	pageNumber := StrToInt(r.URL.Query().Get("page_number"))
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")
	chapterID := r.URL.Query().Get("chapter_id")

	var res = []fileupload.FileUpload{}
	count, err := fileupload.GetFileUploadByChapterID(chapterID, pageSize, pageNumber, sortBy, sortOrder, &res)

	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.SendDataSuccess(w, map[string]interface{}{
			"fileuploads": res,
			"count":       count,
		})
	}
}

func (s *FileUploadServer) HandleGetByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	file, err := fileupload.GetByName(name)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.SendDataSuccess(w, file)
	}
}

func (s *FileUploadServer) mustGetFileupload(r *http.Request) (*fileupload.FileUpload, error) {
	var id = r.URL.Query().Get("id")
	var v, err = fileupload.GetByID(id)
	if err != nil {
		return v, err
	} else {
		return v, nil
	}
}

func (s *FileUploadServer) HandleUpdateFile(w http.ResponseWriter, r *http.Request) {
	var newFileUpload = &fileupload.FileUpload{}
	s.MustDecodeBody(r, newFileUpload)
	v, err := s.mustGetFileupload(r)
	if err != nil {
		s.ErrorMessage(w, "file_not_found")
		return
	}
	err = v.Update(newFileUpload)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		result, err := fileupload.GetByID(v.ID)
		if err != nil {
			s.ErrorMessage(w, "file_not_found")
			return
		}
		s.SendDataSuccess(w, result)
	}
}

func StrToInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}
