package handlers

import (
	"forum/internal/metrics"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"path/filepath"
)

func (h *Handler) Routes() http.Handler {

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("forum"),
		newrelic.ConfigLicense("03be4fa2a486569d91921a396e8efe0fFFFFNRAL"),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		log.Println("cannot create newrelic instance")
	}
	mux := http.NewServeMux()
	// add a css file to route
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.Handle("/metrics", promhttp.Handler())
	metrics.Init()
	mux.HandleFunc(newrelic.WrapHandleFunc(app, "/", h.metricsMiddleware(http.HandlerFunc(h.home))))
	mux.HandleFunc(newrelic.WrapHandleFunc(app, "/login", h.metricsMiddleware(http.HandlerFunc(h.login))))
	mux.HandleFunc(newrelic.WrapHandleFunc(app, "/register", h.metricsMiddleware(http.HandlerFunc(h.register))))
	mux.Handle("/logout", h.requireAuthentication(h.metricsMiddleware(http.HandlerFunc(h.logout))))

	mux.HandleFunc(newrelic.WrapHandleFunc(app, "/post/", h.showPost))
	mux.HandleFunc("/posts", h.GetPosts)
	mux.Handle("/lp", h.requireAuthentication(h.metricsMiddleware(http.HandlerFunc(h.GetLikedPosts))))

	mux.HandleFunc("/postscat", h.showPostsByCategory)
	mux.HandleFunc("/pc", h.GetPostsCat)

	mux.Handle("/myposts", h.requireAuthentication(h.metricsMiddleware(http.HandlerFunc(h.myposts))))
	mux.Handle("/mp", h.requireAuthentication(h.metricsMiddleware(http.HandlerFunc(h.GetMyPosts))))

	mux.Handle("/post/create", h.requireAuthentication(h.metricsMiddleware(http.HandlerFunc(h.createPost))))
	mux.Handle("/post/reaction", h.requireAuthentication(h.metricsMiddleware(http.HandlerFunc(h.reactionPost))))
	mux.Handle("/likedposts", h.requireAuthentication(h.metricsMiddleware(http.HandlerFunc(h.likedPosts))))

	mux.Handle("/comment/create", h.requireAuthentication(h.metricsMiddleware(http.HandlerFunc(h.createComment))))
	mux.Handle("/comment/reaction", h.requireAuthentication(h.metricsMiddleware(http.HandlerFunc(h.reactionComment))))

	return h.authenticate(mux)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
