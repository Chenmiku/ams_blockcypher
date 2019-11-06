package article

import (
	"http/web"
	"myproject/dapi/api/auth/session"
	"myproject/dapi/o/org/article"
	"net/http"
	"strconv"
)

type ArticleServer struct {
	web.JsonServer
	*http.ServeMux
}

func NewArticleServer() *ArticleServer {
	var s = &ArticleServer{
		ServeMux: http.NewServeMux(),
	}

	s.HandleFunc("/create", s.HandleCreate)
	s.HandleFunc("/get", s.HandleGetByID)
	s.HandleFunc("/update", s.HandleUpdateByID)
	s.HandleFunc("/mark_delete", s.HandleMarkDelete)
	s.HandleFunc("/get_all", s.HandleGetAll)
	s.HandleFunc("/get_by_category", s.HandleGetByCategoryID)
	return s
}

func (s *ArticleServer) Session(w http.ResponseWriter, r *http.Request) *session.Session {

	ses, err := session.FromContext(r.Context())

	if err != nil {
		s.SendError(w, err)
	}

	return ses
}

func StrToInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}

func (s *ArticleServer) HandleGetAll(w http.ResponseWriter, r *http.Request) {

	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	pageSize := StrToInt(r.URL.Query().Get("page_size"))
	pageNumber := StrToInt(r.URL.Query().Get("page_number"))

	var res = []article.Article{}

	count, err := article.GetAll(pageSize, pageNumber, sortBy, sortOrder, &res)

	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.SendDataSuccess(w, map[string]interface{}{
			"articles": res,
			"count":    count,
		})
	}
}

func (s *ArticleServer) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var v = &article.Article{}
	s.MustDecodeBody(r, v)
	err := v.Create()
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.SendDataSuccess(w, v)
	}
}

func (s *ArticleServer) mustGetArticle(r *http.Request) (*article.Article, error) {
	var id = r.URL.Query().Get("id")
	var v, err = article.GetByID(id)
	if err != nil {
		return v, err
	} else {
		return v, nil
	}
}

func (s *ArticleServer) HandleUpdateByID(w http.ResponseWriter, r *http.Request) {
	var newArticle = &article.Article{}
	s.MustDecodeBody(r, newArticle)
	v, err := s.mustGetArticle(r)
	if err != nil {
		s.ErrorMessage(w, "article_not_found")
		return
	}
	err = v.Update(newArticle)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		result, err := article.GetByID(v.ID)
		if err != nil {
			s.ErrorMessage(w, "article_not_found")
			return
		}
		s.SendDataSuccess(w, result)
	}
}

func (s *ArticleServer) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	v, err := s.mustGetArticle(r)
	if err != nil {
		s.ErrorMessage(w, "article_not_found")
		return
	}
	s.SendDataSuccess(w, v)
}

func (s *ArticleServer) HandleMarkDelete(w http.ResponseWriter, r *http.Request) {
	v, err := s.mustGetArticle(r)
	if err != nil {
		s.ErrorMessage(w, "article_not_found")
		return
	}
	err = article.MarkDelete(v.ID)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.Success(w)
	}
}

func (s *ArticleServer) HandleGetByCategoryID(w http.ResponseWriter, r *http.Request) {

	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	pageSize := StrToInt(r.URL.Query().Get("page_size"))
	pageNumber := StrToInt(r.URL.Query().Get("page_number"))

	id := r.URL.Query().Get("category_id")

	var res = []article.Article{}

	count, err := article.GetArticleByCategoryID(id, pageSize, pageNumber, sortBy, sortOrder, &res)

	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.SendDataSuccess(w, map[string]interface{}{
			"articles": res,
			"count":    count,
		})
	}
}
