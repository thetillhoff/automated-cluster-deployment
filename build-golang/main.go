package main

// used go version: 1.14

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"text/template"

	"gopkg.in/yaml.v3"
)

// type Values: Created for easier templating
type Values struct {
	Values map[string]interface{}
}

// error checking is often needed, so simplified with
//   check(err)
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Load file and return it's content.
func loadFile(filepath string) string {
	fileContent, err := ioutil.ReadFile(filepath) // read contents from path into variable
	check(err)
	return string(fileContent)
}

// Open file and overwrite it's content.
func writeFile(filepath string, content string) {
	// write the whole body at once
	err := ioutil.WriteFile(filepath, []byte(content), 0644)
	check(err)
}

// Extract IP out of http.Request and return it.
func getUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	u, err := url.Parse("//" + IPAddress)
	check(err)

	return u.Hostname()
}

// Search m recursively for the key 'matchkey' and it's value against 'match'.
// If match is found, return value of key 'searchkey'
// Return found value. Lower level overrides upper level.
func findKey(m map[string]interface{}, matchkey string, match interface{}, searchkey string) (interface{}, error) {

	for key, value := range m { // for each key-value of m
		if mappedValue, ok := value.(map[string]interface{}); ok { // if value is a map
			searchValue, err := findKey(mappedValue, matchkey, match, searchkey) // recursion
			if searchValue != nil && err == nil {
				return searchValue, nil // passthrough value & success
			}
			// matchfound = matchfound wouldn't make sense, right?
		} else { // if value is not a map
			searchValue, ok := m[searchkey]
			if key == matchkey && value == match && ok { // if key equals matchkey, value equals match and searchkey is contained
				return searchValue, nil // return value & success
			}
		}
	}

	return nil, errors.New("findKey: No match found.") // return nil & error
}

func changeKey(m map[string]interface{}, matchkey string, match interface{}, changekey string, changevalue interface{}) error {

	for key, value := range m { // for each key-value of m
		if mappedValue, ok := value.(map[string]interface{}); ok { // if value is a map
			err := changeKey(mappedValue, matchkey, match, changekey, changevalue) // recursion
			if err == nil {
				return nil // passthrough success
			}
		} else { // if value is not a map
			_, ok := m[changekey] // check if key 'changekey' is contained
			if key == matchkey && value == match && ok {
				m[changekey] = changevalue
				return nil // return success
			}
		}
	}

	return errors.New("changeKey: No match found.") // return error

}

// Search m recursively for the key 'matchkey' and it's value against 'match'.
// If found, return it's parent key.
func findParentKey(m map[string]interface{}, matchkey string, match interface{}) (interface{}, error) {

	for key, value := range m { // for each key-value of m
		if mappedValue, ok := value.(map[string]interface{}); ok { // if value is a map
			parentKey, err := findParentKey(mappedValue, matchkey, match) // recursion
			if parentKey != nil && err == nil {
				return parentKey, nil // passthrough key & success
			} else if parentKey == nil && err == nil { // match found, but parent key empty
				return key, nil // return key & success
			}
		} else { // if value is no map
			if key == matchkey && value == match {
				return nil, nil // return match found, so take the parent
			}
		}
	}

	return nil, errors.New("findParentKey: No match or no parent found.") // return nil & error
}

// check if request container get parameter exactly once and return it (and true if success | false for error)
func GETsingle(r *http.Request, get string) (string, bool) {
	gets, ok := r.URL.Query()["mac"] // could be multiple 'mac' parameters, thus error checking follows
	if !ok || len(gets) == 0 {       // if parameter 'mac' was not found
		return "No " + get + " provided.", false
	} else if ok && len(gets) > 1 { // if more than one parameter 'mac' was found
		return "More than one " + get + " provided.", false
	}

	// exactly one mac was provided
	return gets[0], true // extract the single mac
}

// validate 'mac' for being a valid MAC
func validMAC(mac string) bool {
	re := regexp.MustCompile("^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$")
	return re.Match([]byte(mac))
}

// Handler for '<SERVERADDR>/default'
func defaultHandler(w http.ResponseWriter, r *http.Request) {

	mac, ok := GETsingle(r, "mac")
	if !ok {
		fmt.Fprintf(w, mac)
	} else if !validMAC(mac) {
		fmt.Fprintf(w, "Provided MAC is invalid.")
	} else { // mac was provided correctly
		content := loadFile("/etc/ansible/hosts")             // load file contents
		mappedContent := make(map[string]interface{}, 0)      // create empty map
		yaml.Unmarshal([]byte(content), &mappedContent)       // store yaml in map
		value, _ := findKey(mappedContent, "MAC", mac, "MAC") // search for host with correct mac // don't do error tracking, as 'not found' error isn't important here
		if value == nil {                                     // if host not found
			newhost := make(map[string]interface{})
			// setting default parameters for unkown hosts
			newhost["MAC"] = mac
			newhost["managed"] = false
			newhost["IP"] = getUserIP(r)
			newhost["state"] = "none"

			all, ok := mappedContent["all"] // check if mappedContent contains key 'all'
			if !ok || all == nil {          // if not
				all = make(map[string]interface{}, 0) // create it
				mappedContent["all"] = all            // and add it to mappedContent
			}
			hosts, ok := all.(map[string]interface{})["hosts"] // check if all container key 'hosts'
			if !ok || hosts == nil {                           // if not
				hosts = make(map[string]interface{}, 0)       // create it
				all.(map[string]interface{})["hosts"] = hosts // and add it to all
			}
			if hosts, ok := hosts.(map[string]interface{}); ok { // if 'hosts' is map (which it should be...)
				hosts[mac] = newhost

				newcontent, err := yaml.Marshal(&mappedContent) // store map into yaml
				check(err)
				writeFile("/etc/ansible/hosts", string(newcontent))

				t, err := template.ParseFiles("ipxe_menu" + ".tmpl")
				check(err)
				m := make(map[string]interface{})
				m["message"] = "Added this host to the inventory as ungrouped. If managed, reboot after state is set."
				data := Values{
					Values: m,
				}
				err = t.Execute(w, data)
				check(err)
			} else {
				// 'hosts' isn't a map
				fmt.Fprintf(w, "all.hosts isn't a map.")
			}
		} else { // else: host is found
			couldbemanaged, err := findKey(mappedContent, "MAC", mac, "managed") // get state of host
			check(err)
			managed := couldbemanaged.(bool)
			if managed { // if host is managed

				couldbestate, err := findKey(mappedContent, "MAC", mac, "state") // get state of host
				check(err)
				state := couldbestate.(string)
				switch state {
				case "waiting for provisioning":
					t, err := template.ParseFiles("provisioning_debian" + ".tmpl")
					check(err)
					m := make(map[string]interface{})

					// server
					m["server"] = r.Host // get own hostname

					// hostname
					hostname, err := findParentKey(mappedContent, "MAC", mac) // get parent of MAC-entry
					check(err)
					m["hostname"] = hostname

					// execute template
					data := Values{
						Values: m,
					}
					err = t.Execute(w, data)
					check(err)
				case "offline":
					t, err := template.ParseFiles("ipxe_localboot" + ".tmpl")
					check(err)
					m := make(map[string]interface{})

					// server
					m["server"] = r.Host // get own hostname

					// hostname
					hostname, err := findParentKey(mappedContent, "MAC", mac) // get parent of MAC-entry
					check(err)
					m["hostname"] = hostname

					// execute template
					data := Values{
						Values: m,
					}
					err = t.Execute(w, data)
					check(err)

					changeKey(mappedContent, "MAC", mac, "state", "booting from local device") // change state to provisioning
					newcontent, err := yaml.Marshal(&mappedContent)                            // store map into yaml
					check(err)
					writeFile("/etc/ansible/hosts", string(newcontent))
				default: // e.g. state == 'none' or state == nil
					t, err := template.ParseFiles("ipxe_menu" + ".tmpl")
					check(err)
					m := make(map[string]interface{})
					m["message"] = "No or wrong state set for this host. Expected state was either '" + "waiting for provisioning" + "' or '" + "offline" + "', actual state was '" + state + "'."
					data := Values{
						Values: m,
					}
					err = t.Execute(w, data)
					check(err)
				}
			} else { // host is unmanaged
				// display menu

				t, err := template.ParseFiles("ipxe_menu" + ".tmpl")
				check(err)
				m := make(map[string]interface{})
				m["message"] = "This host is unmanaged. Continuing with menu."
				data := Values{
					Values: m,
				}
				err = t.Execute(w, data)
				check(err)
			}
		}
	}
}

// Handler for '<SERVERADDR>/preseed'
func preseedHandler(w http.ResponseWriter, r *http.Request) {

	mac, ok := GETsingle(r, "mac")
	if !ok {
		fmt.Fprintf(w, mac)
	} else if !validMAC(mac) {
		fmt.Fprintf(w, "Provided MAC is invalid.")
	} else { // mac was provided correctly
		content := loadFile("/etc/ansible/hosts")             // load file contents
		mappedContent := make(map[string]interface{}, 0)      // create empty map
		yaml.Unmarshal([]byte(content), &mappedContent)       // store yaml in map
		value, _ := findKey(mappedContent, "MAC", mac, "MAC") // search for host with correct mac // don't do error tracking, as 'not found' error isn't important here
		if value == nil {                                     // if host not found
			fmt.Fprintf(w, "No host with provided MAC.")
		} else { // else: host is found
			couldbemanaged, err := findKey(mappedContent, "MAC", mac, "managed") // get state of host
			check(err)
			managed := couldbemanaged.(bool)
			if managed { // if host is managed

				couldbestate, err := findKey(mappedContent, "MAC", mac, "state") // get state of host
				check(err)
				state := couldbestate.(string)
				switch state {
				case "waiting for provisioning":
					// -> correct state
					t, err := template.ParseFiles("preseed" + ".tmpl")
					check(err)
					m := make(map[string]interface{})

					// server
					m["server"] = r.Host // get own hostname

					// mac
					m["mac"] = mac

					// hostname
					hostname, err := findParentKey(mappedContent, "MAC", mac) // get parent of MAC-entry
					check(err)
					m["hostname"] = hostname

					// username
					m["username"] = "enforge"

					// pass
					m["pass"] = "yftK48L59TcL6" // hash of somepass

					// mirror
					m["mirror"] = "deb.debian.org"

					// packages
					m["packages"] = "openssh-server wget curl git net-tools nano"

					// execute template
					data := Values{
						Values: m,
					}
					err = t.Execute(w, data)
					check(err)

					changeKey(mappedContent, "MAC", mac, "state", "provisioning") // change state to provisioning
					newcontent, err := yaml.Marshal(&mappedContent)               // store map into yaml
					check(err)
					writeFile("/etc/ansible/hosts", string(newcontent))
				default: // e.g. state == 'none' or state == nil
					t, err := template.ParseFiles("ipxe_menu" + ".tmpl")
					check(err)
					m := make(map[string]interface{})
					m["message"] = "No or wrong state set for this host. Expected state was '" + "waiting for provisioning" + "', actual state was '" + state + "'."
					data := Values{
						Values: m,
					}
					err = t.Execute(w, data)
					check(err)
				}
			} else { // host is unmanaged
				// display menu

				t, err := template.ParseFiles("ipxe_menu" + ".tmpl")
				check(err)
				m := make(map[string]interface{})
				m["message"] = "This host is unmanaged. Continuing with menu."
				data := Values{
					Values: m,
				}
				err = t.Execute(w, data)
				check(err)
			}
		}
	}
}

// Handler for '<SERVERADDR>/preseedlate'
func preseedlateHandler(w http.ResponseWriter, r *http.Request) {

	mac, ok := GETsingle(r, "mac")
	if !ok {
		fmt.Fprintf(w, mac)
	} else if !validMAC(mac) {
		fmt.Fprintf(w, "Provided MAC is invalid.")
	} else { // mac was provided correctly
		content := loadFile("/etc/ansible/hosts")             // load file contents
		mappedContent := make(map[string]interface{}, 0)      // create empty map
		yaml.Unmarshal([]byte(content), &mappedContent)       // store yaml in map
		value, _ := findKey(mappedContent, "MAC", mac, "MAC") // search for host with correct mac // don't do error tracking, as 'not found' error isn't important here
		if value == nil {                                     // if host not found
			fmt.Fprintf(w, "No host with provided MAC.")
		} else { // else: host is found
			couldbemanaged, err := findKey(mappedContent, "MAC", mac, "managed") // get state of host
			check(err)
			managed := couldbemanaged.(bool)
			if managed { // if host is managed

				couldbestate, err := findKey(mappedContent, "MAC", mac, "state") // get state of host
				check(err)
				state := couldbestate.(string)
				switch state {
				case "provisioning":

					t, err := template.ParseFiles("preseedlate" + ".tmpl")
					check(err)
					m := make(map[string]interface{})

					// server
					m["server"] = r.Host // get own hostname

					// mac
					m["mac"] = mac

					// username
					m["username"] = "enforge"

					// execute template
					data := Values{
						Values: m,
					}
					err = t.Execute(w, data)
					check(err)

					changeKey(mappedContent, "MAC", mac, "state", "booting from local device") // change state to offline
					newcontent, err := yaml.Marshal(&mappedContent)                            // store map into yaml
					check(err)
					writeFile("/etc/ansible/hosts", string(newcontent))
				default: // e.g. state == 'none' or state == nil
					t, err := template.ParseFiles("ipxe_menu" + ".tmpl")
					check(err)
					m := make(map[string]interface{})
					m["message"] = "No or wrong state set for this host. Expected state was '" + "provisioning" + "', actual state was '" + state + "'."
					data := Values{
						Values: m,
					}
					err = t.Execute(w, data)
					check(err)
				}
			} else { // host is unmanaged
				// display menu

				t, err := template.ParseFiles("ipxe_menu" + ".tmpl")
				check(err)
				m := make(map[string]interface{})
				m["message"] = "This host is unmanaged. Continuing with menu."
				data := Values{
					Values: m,
				}
				err = t.Execute(w, data)
				check(err)
			}
		}
	}
}

// Handler for '<SERVERADDR>/hostonline'
func hostonlineHandler(w http.ResponseWriter, r *http.Request) {

	mac, ok := GETsingle(r, "mac")
	if !ok {
		fmt.Fprintf(w, mac)
	} else if !validMAC(mac) {
		fmt.Fprintf(w, "Provided MAC is invalid.")
	} else { // mac was provided correctly
		content := loadFile("/etc/ansible/hosts")             // load file contents
		mappedContent := make(map[string]interface{}, 0)      // create empty map
		yaml.Unmarshal([]byte(content), &mappedContent)       // store yaml in map
		value, _ := findKey(mappedContent, "MAC", mac, "MAC") // search for host with correct mac // don't do error tracking, as 'not found' error isn't important here
		if value == nil {                                     // if host not found
			fmt.Fprintf(w, "No host with provided MAC.")
		} else { // else: host is found
			couldbemanaged, err := findKey(mappedContent, "MAC", mac, "managed") // get state of host
			check(err)
			managed := couldbemanaged.(bool)
			if managed { // if host is managed

				couldbestate, err := findKey(mappedContent, "MAC", mac, "state") // get state of host
				check(err)
				state := couldbestate.(string)
				switch state {
				case "booting from local device":

					err := changeKey(mappedContent, "MAC", mac, "state", "online") // change state to provisioning
					check(err)

					fmt.Fprintln(w, "Welcome! Changed 'state=online' for your host set in the inventory.")

					ip := getUserIP(r)
					existingip, err := findKey(mappedContent, "MAC", mac, "IP")
					check(err)
					if ip != existingip.(string) {
						err := changeKey(mappedContent, "MAC", mac, "IP", ip)
						check(err)
						fmt.Fprintln(w, "Your IP address changed from "+existingip.(string)+" to "+ip+".")
					}

					newcontent, err := yaml.Marshal(&mappedContent) // store map into yaml
					check(err)
					writeFile("/etc/ansible/hosts", string(newcontent))
				default: // e.g. state == 'none' or state == nil
					t, err := template.ParseFiles("ipxe_menu" + ".tmpl")
					check(err)
					m := make(map[string]interface{})
					m["message"] = "No or wrong state set for this host. Expected state was '" + "booting from local device" + "', actual state was '" + state + "'."
					data := Values{
						Values: m,
					}
					err = t.Execute(w, data)
					check(err)
				}
			} else { // host is unmanaged
				// display menu

				t, err := template.ParseFiles("ipxe_menu" + ".tmpl")
				check(err)
				m := make(map[string]interface{})
				m["message"] = "This host is unmanaged. Continuing with menu."
				data := Values{
					Values: m,
				}
				err = t.Execute(w, data)
				check(err)
			}
		}
	}
}

// Handler for '<SERVERADDR>/healthcheck'
func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "up")
}

func main() {
	fmt.Print("To access this webserver access localhost:8080\n")
	http.HandleFunc("/healthcheck", healthcheckHandler)
	http.HandleFunc("/default", defaultHandler)
	http.HandleFunc("/preseed", preseedHandler)
	http.HandleFunc("/preseedlate", preseedlateHandler)
	http.HandleFunc("/hostonline", hostonlineHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
