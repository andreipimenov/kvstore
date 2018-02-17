package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/andreipimenov/kvstore/model"
)

//Client - client implementation for interacting with "kvstore" server
type Client struct {
	ServerHost      string
	ServerPort      int
	Login           string
	Password        string
	AuthorizedToken string
	Client          *http.Client
}

//NewClient creates client
func NewClient(serverHost string, serverPort int, login string, password string) *Client {
	return &Client{
		ServerHost: serverHost,
		ServerPort: serverPort,
		Login:      login,
		Password:   password,
		Client: &http.Client{
			Timeout: time.Duration(30 * time.Second),
		},
	}
}

//ProcessRequest implements request to api and returns string representation of response OR error
func (c *Client) ProcessRequest(method string, uri string, body io.Reader) string {
	req, err := http.NewRequest(method, fmt.Sprintf("http://%s:%d/api/v1%s", c.ServerHost, c.ServerPort, uri), body)
	if err != nil {
		return fmt.Sprintf("ERROR: %s\n", err.Error())
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Sprintf("ERROR: %s\n", err.Error())
	}
	if resp == nil {
		return fmt.Sprintf("ERROR: nil response\n")
	}
	defer resp.Body.Close()
	v, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("ERROR: %s\n", err.Error())
	}
	return string(v)
}

//Ping - healthcheck
func (c *Client) Ping() string {
	return c.ProcessRequest(http.MethodGet, "/ping", nil)
}

//Get returns value by key or error representation
func (c *Client) Get(key string) string {
	return c.ProcessRequest(http.MethodGet, fmt.Sprintf("/keys/%s/values", key), nil)
}

//GetIndex returns indexed value from list or map
func (c *Client) GetIndex(key string, index interface{}) string {
	return c.ProcessRequest(http.MethodGet, fmt.Sprintf("/keys/%s/values/%v", key, index), nil)
}

//Set - set value by key
func (c *Client) Set(key string, value string) string {
	var v interface{}
	err := json.Unmarshal([]byte(value), &v)
	if err != nil {
		//Here the big limitation cause string respresentation of json will being escaped
		//TODO: check if data is string OR json-style list/map and make if flexible
	 	v = value
	}
	j, _ := json.Marshal(&model.APIKeyValue{
		Key: key, Value: v,
	})
	return c.ProcessRequest(http.MethodPost, "/keys", bytes.NewReader(j))
}

//Remove - remove key
func (c *Client) Remove(key string) string {
	return c.ProcessRequest(http.MethodDelete, fmt.Sprintf("/keys/%s", key), nil)
}

//Keys - returns keys by pattern
func (c *Client) Keys(pattern string) string {
	return c.ProcessRequest(http.MethodGet, fmt.Sprintf("/keys/%s", pattern), nil)
}

//SetExpires - set expiration time for key
func (c *Client) SetExpires(key string, expires int64) string {
	j, _ := json.Marshal(&model.APIKeyExpires{
		Expires: expires,
	})
	return c.ProcessRequest(http.MethodPost, fmt.Sprintf("/keys/%s/expires", key), bytes.NewReader(j))
}

//GetExpires - returns expiration time for key
func (c *Client) GetExpires(key string) string {
	return c.ProcessRequest(http.MethodGet, fmt.Sprintf("/keys/%s/expires", key), nil)
}

//WebUI - simple web user interface
func (c *Client) WebUI() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
			<!DOCTYPE html>
			<html>
			<head>
				<meta charset="utf-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=0">
				<title>KVStore WEB Client</title>
				<link href="https://fonts.googleapis.com/css?family=Roboto+Mono&amp;subset=cyrillic" rel="stylesheet">
				<style type="text/css">
				* {
					margin: 0;
					padding: 0;
					box-sizing: border-box;
				}
				html, body {
					widht: 100%;
					height: 100%;
				}
				body {
					font-size: 14px;
					line-height: 24px;
					color: #eee;
					background: #2c3e50;
					font-family: 'Roboto Mono';
					display: flex;
					flex-wrap: wrap;
					justify-content: center;
					align-content: center;
				}
				body > div {
					width: 90%;
					padding: 10px;
				}
				@media (min-width: 768px) {
				body > div {
					width: 70%;
				}
				}
				div.form {
					padding: 10px 0;
					display: flex;
					flex-wrap: wrap;
				}
				div.form > div {
					margin-bottom: 12px;
					width: 100%;
					display: flex;
					justify-content: flex-start;
					flex-wrap: wrap;
				}
				div.form > div > p {
					width: 100%;
				}
				input, select, textarea {
					background: #566b7f;
					width: 100px;
					max-width: 25%;
					border: 0;
					padding: 10px 15px;
					color: #fff;
					border-radius: 2px;
					font-family: 'Roboto Mono';
					margin-right: 2px;
				}
				input, select {
					width: 100%;
				}
				textarea {
					width: 200px;
					min-height: 48px;
				}
				button {
					display: inline-block;
					background: #16a085;
					border: 0;
					color: #eee;
					max-width: 50%;
					padding: 10px 15px;
					border-radius: 2px;
					cursor: pointer;
					font-family: 'Roboto Mono';
				}
				</style>
			</head>
			<body>
				<div class="request">
					<div class="form">
						<div>
							<select name="command">
								<option value="ping">PING</option>
								<option value="set">SET</option>
								<option value="get">GET</option>
								<option value="getindex">GET INDEX</option>
								<option value="remove">REMOVE</option>
								<option value="keys">KEYS</option>
								<option value="setexpires">SET EXPIRES</option>
								<option value="getexpires">GET EXPIRES</option>
							</select>
							<input type="text" name="first">
							<textarea type="text" name="second"></textarea>
						</div>
						<button onclick="process()">Go</button>
					</div>
				</div>
				<div class="response">
					<p id="response"></p>
				</div>
			<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.12.2/jquery.min.js"></script>
			<script>
			$(document).keypress(function(e) {
				var keycode = (e.keyCode ? e.keyCode : e.which);
				if (keycode == '13') {
					process();
				}
			});
			$('select').change(function() {
				$('input, textarea').val('');
				if ($(this).val() == 'keys') {
					$('input[name=first]').val('*');
				}
			});
			function process() {
				var data = [];
				$('input, select, textarea').each(function(i, elem) {
					var name = $(elem).attr('name');
					var val = $(elem).val();
					data.push($(elem).attr('name')+'='+encodeURIComponent(val));
				});
				var body = data.join('&');
				console.log(body);
				$.ajax({
					url: '/process',
					type: 'POST',
					data: body,
					dataType: 'json'
				})
				.done(function(data, textStatus, jqXHR) {
					$('#response').text(jqXHR.responseText);
				})
				.fail(function(jqXHR, textStatus, errorThrown) {
					$('#response').text(jqXHR.responseText);
				});
			}
			</script>
			</body>
			</html>
		`)
	})
}

//ProcessWebUI - send request taken from WebUI to KVServer and return response back to browser
func (c *Client) ProcessWebUI() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		command := r.PostFormValue("command")
		first := r.PostFormValue("first")
		second := r.PostFormValue("second")
		switch command {
		case "ping":
			fmt.Fprint(w, c.Ping())
		case "set":
			fmt.Fprint(w, c.Set(first, second))
		case "get":
			fmt.Fprint(w, c.Get(first))
		case "getindex":
			fmt.Fprint(w, c.GetIndex(first, second))
		case "remove":
			fmt.Fprint(w, c.Remove(first))
		case "keys":
			fmt.Fprint(w, c.Keys(first))
		case "setexpires":
			expires, _ := strconv.ParseInt(second, 10, 64)
			fmt.Fprint(w, c.SetExpires(first, expires))
		case "getexpires":
			fmt.Fprint(w, c.GetExpires(first))
		default:
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}
	})
}
