// go:generate template_compiler
package main

import "github.com/googleapis/openapi-compiler/plugins/go/template_plugin"

func main() {
    template_plugin.Run(map[string]string{ 
        "client.go": "cGFja2FnZSB7ey5SZW5kZXJlci5QYWNrYWdlfX0KCmltcG9ydCAoCiAgImJ5dGVzIgogICJlbmNvZGluZy9qc29uIgogICJmbXQiCiAgIm5ldC9odHRwIgogICJzdHJpbmdzIgopCiAgCmNvbnN0IHNlcnZpY2UgPSAiaHR0cDovL2xvY2FsaG9zdDo4MDgwIgoKdHlwZSBDbGllbnQgc3RydWN0IHsKCVNlcnZpY2Ugc3RyaW5nCn0gCgpmdW5jIE5ld0NsaWVudChzZXJ2aWNlIHN0cmluZykgKkNsaWVudCB7CgljbGllbnQgOj0gJkNsaWVudHt9CgljbGllbnQuU2VydmljZSA9IHNlcnZpY2UKCXJldHVybiBjbGllbnQKfQoKe3tyYW5nZSAuUmVuZGVyZXIuTWV0aG9kc319Cnt7aWYgZXEgLlJlc3VsdFR5cGVOYW1lICIifX0KZnVuYyAoY2xpZW50ICpDbGllbnQpIHt7LkNsaWVudE5hbWV9fSh7e3BhcmFtZXRlckxpc3QgLn19KSAoZXJyIGVycm9yKSB7Cnt7ZWxzZX19CmZ1bmMgKGNsaWVudCAqQ2xpZW50KSB7ey5DbGllbnROYW1lfX0oe3twYXJhbWV0ZXJMaXN0IC59fSkgKHJlc3VsdCAqe3suUmVzdWx0VHlwZU5hbWV9fSwgZXJyIGVycm9yKSB7Cnt7ZW5kfX0KCXBhdGggOj0gY2xpZW50LlNlcnZpY2UgKyAie3suUGF0aH19IgoJCgl7e3JhbmdlIC5QYXJhbWV0ZXJzVHlwZS5GaWVsZHN9fQkKCXt7aWYgZXEgLlBvc2l0aW9uICJwYXRoIn19CglwYXRoID0gc3RyaW5ncy5SZXBsYWNlKHBhdGgsICJ7IiArICJ7ey5KU09OTmFtZX19IiArICJ9IiwgZm10LlNwcmludGYoIiV2Iiwge3suSlNPTk5hbWV9fSksIDEpCgl7e2VuZH19Cgl7e2VuZH19CgkKCXt7aWYgZXEgLk1ldGhvZCAiUE9TVCJ9fQoJYm9keSA6PSBuZXcoYnl0ZXMuQnVmZmVyKQoJanNvbi5OZXdFbmNvZGVyKGJvZHkpLkVuY29kZSh7e2JvZHlQYXJhbWV0ZXJOYW1lIC59fSkKCXJlcSwgZXJyIDo9IGh0dHAuTmV3UmVxdWVzdCgie3suTWV0aG9kfX0iLCBwYXRoLCBib2R5KQoJe3tlbHNlfX0KCXJlcSwgZXJyIDo9IGh0dHAuTmV3UmVxdWVzdCgie3suTWV0aG9kfX0iLCBwYXRoLCBuaWwpCgl7e2VuZH19CgkKCWlmIGVyciAhPSBuaWwgewoJCXJldHVybgoJfQoJcmVzcCwgZXJyIDo9IGh0dHAuRGVmYXVsdENsaWVudC5EbyhyZXEpCglpZiBlcnIgIT0gbmlsIHsKCQlyZXR1cm4gCgl9CglkZWZlciByZXNwLkJvZHkuQ2xvc2UoKQoJe3tpZiBuZSAuUmVzdWx0VHlwZU5hbWUgIiJ9fQoJZGVjb2RlciA6PSBqc29uLk5ld0RlY29kZXIocmVzcC5Cb2R5KQoJcmVzdWx0ID0gJnt7LlJlc3VsdFR5cGVOYW1lfX17fQoJZGVjb2Rlci5EZWNvZGUocmVzdWx0KQoJe3tlbmR9fQoJcmV0dXJuICAKfQp7e2VuZH19Cg==",
        "types.go": "cGFja2FnZSB7ey5SZW5kZXJlci5QYWNrYWdlfX0KCmltcG9ydCAoCikKCi8vIHR5cGVzCnt7cmFuZ2UgLlJlbmRlcmVyLlR5cGVzfX0KdHlwZSB7ey5OYW1lfX0gc3RydWN0IHsge3tyYW5nZSAuRmllbGRzfX0KICB7ey5OYW1lfX0ge3suVHlwZX19e3tpZiBuZSAuSlNPTk5hbWUgIiJ9fSBganNvbjoie3suSlNPTk5hbWV9fSJge3tlbmR9fXt7ZW5kfX0KfQp7e2VuZH19Cg==",
    })
}