package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	rawurl := r.URL.String()[1:]
	u, err := url.Parse(rawurl)
	if checkErr(err) {
		return
	}
	isMainReq := u.Scheme != ""
	if !isMainReq { // subsequent requests
		var referer string
		cookie, err := r.Cookie("proxlet-host")
		if err == nil {
			referer = cookie.Value
		} else {
			// Fallback to Referer header, which does not always work.
			tmp, err := url.Parse(r.Header.Get("Referer"))
			if checkErr(err) {
				return
			}
			referer = tmp.Path[1:]
			r.Header.Set("Referer", referer)
		}
		ru, err := url.Parse(referer)
		if checkErr(err) {
			return
		}
		u, err = url.Parse(ru.Scheme + "://" + ru.Host + "/" + rawurl)
		if checkErr(err) {
			return
		}
	}
	//log.Println(u.String())
	r.URL = u
	r.Host = u.Host
	if rcookie, ok := r.Header["Cookie"]; ok {
		for i, v := range rcookie {
			if strings.HasPrefix(v, "proxlet-host=") {
				rcookie = append(rcookie[:i], rcookie[i+1:]...)
				break
			}
		}
	}
	r.Header.Set("Accept-Encoding", "identity") // disable any encoding
	h := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
	})
	h.ModifyResponse = func(rsp *http.Response) (err error) { // server-side redirect
		if isMainReq {
			defer func() {
				http.SetCookie(w, &http.Cookie{
					Name:  "proxlet-host",
					Value: r.URL.String(),
					Path:  "/",
				})
			}()
		}
		switch rsp.StatusCode {
		case 301, 302, 303, 307, 308:
		default:
			return nil
		}
		loc, err := r.URL.Parse(rsp.Header.Get("Location"))
		if err != nil {
			return
		}
		//log.Println(loc.String())
		r.URL = loc
		r.Host = r.URL.Host
		r.RequestURI = ""
		rsp_, err := http.DefaultClient.Do(r)
		if rsp_ != nil {
			*rsp = *rsp_
		}
		return
	}
	h.ServeHTTP(w, r)
}

func checkErr(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}
