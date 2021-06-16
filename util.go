package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"time"
)

func checkRunning(host, port string) error {
	attempts := 10
	timeout := 100 * time.Millisecond

	for {
		if attempts == 0 {
			break
		}

		conn, _ := net.DialTimeout("tcp", net.JoinHostPort(host, port), 100*time.Millisecond)
		if conn == nil {
			return nil
		}
		conn.Close()

		time.Sleep(timeout)

		timeout = timeout * 2
		attempts--
	}

	return fmt.Errorf("giving up waiting for %s:%s", host, port)
}

func mustNewTemplate(template *template.Template, err error) *template.Template {
	if err != nil {
		log.Fatal(err)
	}
	return template
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
