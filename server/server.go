package server

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/bayesianmind/demo-file-server/filestore"
	"github.com/bayesianmind/demo-file-server/userstore"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//type Interface interface {
//	RegisterUser(ctx context.Context, user, password string) error
//	Login(ctx context.Context, user, password string) (*LoginResponse, error)
//	PutUserFile(ctx context.Context, user string, file io.Reader) error
//	GetUserFile(ctx context.Context, user, path string) (io.ReadCloser, error)
//	DeleteUserFile(ctx context.Context, user, path string) error
//	ListFiles(ctx context.Context, user string)
//}

type Server struct {
	fstore filestore.Interface
	ustore userstore.Interface
}

func New(fstore filestore.Interface, ustore userstore.Interface) *Server {
	return &Server{fstore: fstore, ustore: ustore}
}

// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (s *Server) Run(addr ...string) error {
	r := s.setupRoutes()
	return r.Run(addr...)
}

type loginT struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type errResponseT struct {
	Err string `json:"error"`
}

func (s *Server) regHandler(c *gin.Context) {
	login := &loginT{}
	err := c.ShouldBindWith(login, binding.JSON)
	if err != nil {
		handleError(c, "invalid json: ", err, 403)
		return
	}
	err = s.ustore.RegisterUser(c, login.Username, login.Password)
	if err != nil {
		handleError(c, "could not register: ", err, 403)
		return
	}
}

func (s *Server) loginHandler(c *gin.Context) {
	login := &loginT{}
	err := c.ShouldBindWith(login, binding.JSON)
	if err != nil {
		handleError(c, "invalid json: ", err, 403)
		return
	}
	resp, err := s.ustore.Login(c, login.Username, login.Password)
	if err != nil {
		handleError(c, "could not log in: ", err, 403)
		return
	}

	c.JSON(200, gin.H{
		"token": resp.SessionToken,
	})
}

func (s *Server) putFileHandler(c *gin.Context) {
	user := c.GetString(userKey) // securely set by middleware
	path := c.Param("fname")
	contentType := c.GetHeader("Content-Type")
	err := s.fstore.PutUserFile(c, user, path, c.Request.Body, contentType)
	if err != nil {
		handleError(c, "error writing file: ", err, 500)
		return
	}
	c.Header("Location", "/files/"+path)
	c.Status(201)
}

func (s *Server) getFileHandler(c *gin.Context) {
	user := c.GetString(userKey) // securely set by middleware
	path := c.Param("fname")
	file, ct, err := s.fstore.GetUserFile(c, user, path)
	if os.IsNotExist(err) {
		c.String(404, "File not found")
		return
	}
	if err != nil {
		handleError(c, "error getting file: ", err, 500)
		return
	}
	if ct != "" { // let Go detect content type if we don't have explicit ctype
		c.Header("Content-Type", ct)
	}
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.String(500, err.Error())
		return
	}
}

func (s *Server) deleteFileHandler(c *gin.Context) {
	user := c.GetString(userKey) // securely set by middleware
	path := c.Param("fname")
	_, _, err := s.fstore.GetUserFile(c, user, path)
	if os.IsNotExist(err) {
		c.String(404, "File not found")
		return
	}
	if err != nil {
		handleError(c, "error stating file: ", err, 500)
		return
	}
	err = s.fstore.DeleteUserFile(c, user, path)
	if err != nil {
		handleError(c, "error deleting file: ", err, 500)
		return
	}
	c.Status(204)
}

func (s *Server) listFilesHandler(c *gin.Context) {
	user := c.GetString(userKey) // securely set by middleware
	files, err := s.fstore.ListFiles(c, user)
	if err != nil {
		handleError(c, "error listing files: ", err, 500)
		return
	}
	c.Header("Content-Type", "application/json")
	if len(files) == 0 {
		// Go allows nil slices of size 0, but they serialize as NULL so handle that case
		_, err := c.Writer.Write([]byte("[]"))
		if err != nil {
			handleError(c, "error writing json: ", err, 500)
		}
		return
	}
	body, err := json.Marshal(files)
	if err != nil {
		handleError(c, "error writing json: ", err, 500)
		return
	}
	c.Writer.Write(body) // 200
}

func handleError(c *gin.Context, msg string, err error, code int) {
	fmt.Println(msg, err)
	c.JSON(code, errResponseT{Err: msg + err.Error()})
}

func (s *Server) setupRoutes() *gin.Engine {
	r := gin.Default()
	r.POST("/login", s.loginHandler)
	r.POST("/register", s.regHandler)

	authorized := r.Group("/", requireValidSession)
	authorized.PUT("/files/:fname", s.putFileHandler)
	authorized.GET("/files/:fname", s.getFileHandler)
	authorized.DELETE("/files/:fname", s.deleteFileHandler)
	authorized.GET("/files", s.listFilesHandler)
	return r
}
