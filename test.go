// package main

// import (
//     "fmt"
//     "strings"
// 	"github.com/openconfig/gnmi/proto/gnmi"
// )

// func main() {
//     parentPath1 := &gnmi.Path{
//         Elem: []*gnmi.PathElem{
//             {Name: "openconfig-interfaces:interfaces"},
//             {Name: "interface", Key: map[string]string{"name": "eth0"}},
//         },
//     }

//     childPath1 := &gnmi.Path{
//         Elem: []*gnmi.PathElem{
//             {Name: "interfaces"},
//             {Name: "interface", Key: map[string]string{"name": "eth0"}},
//             {Name: "state"},
//         },
//     }

//     parentPath2 := &gnmi.Path{
//         Elem: []*gnmi.PathElem{
//             {Name: "openconfig-interfaces:interfaces"},
//             {Name: "interface"},
//         },
//     }

//     childPath2 := &gnmi.Path{
//         Elem: []*gnmi.PathElem{
//             {Name: "interfaces"},
//             {Name: "interface", Key: map[string]string{"name": "eth0"}},
//             {Name: "state"},
//         },
//     }

// 	parentPath3 := &gnmi.Path{
//         Elem: []*gnmi.PathElem{
//             {Name: "openconfig-interfaces:interfaces"},
//             {Name: "interface"},
//         },
//     }

//     childPath3 := &gnmi.Path{
//         Elem: []*gnmi.PathElem{
//             {Name: "interfaces"},
//             {Name: "interface", Key: map[string]string{"name": "eth0"}},
//         },
//     }

//     if pathContains(parentPath1, childPath1) {
//         fmt.Println("Parent path 1 contains the child path 1.")
//     } else {
//         fmt.Println("Parent path 1 does not contain the child path 1.")
//     }

//     if pathContains(parentPath2, childPath2) {
//         fmt.Println("Parent path 2 contains the child path 2.")
//     } else {
//         fmt.Println("Parent path 2 does not contain the child path 2.")
//     }

//     if pathContains(parentPath3, childPath3) {
//         fmt.Println("Parent path 3 contains the child path 3.")
//     } else {
//         fmt.Println("Parent path 3 does not contain the child path 3.")
//     }
// }

package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/openconfig/gnmi/proto/gnmi"
)

func extractModuleName(elem *gnmi.PathElem) string {
    // Split the element name by ":"
    parts := strings.Split(elem.Name, ":")
    if len(parts) > 1 {
        return parts[1]
    }
    return elem.Name
}

func extractWords(input string) string {
    // Define a regular expression pattern to match characters to be removed
    pattern := `[^\w:-]+`
    
    // Compile the regular expression
    regex := regexp.MustCompile(pattern)
    
    // Replace matched characters with an empty string
    result := regex.ReplaceAllString(input, "")
    
    return result
}

func pathContains(parent, child *gnmi.Path) bool {
    // Check if the parent path has fewer elements than the child path
    if len(parent.Elem) > len(child.Elem) {
        return false
    }

    // Iterate over the elements of the parent and child paths
    for i := range parent.Elem {
        withoutModuleNameParent := extractModuleName(child.Elem[i])
        withoutModuleNameChild := extractModuleName(child.Elem[i])

        if withoutModuleNameParent != withoutModuleNameChild {
            return false
        }

        // Check if the keys of the parent element are not a subset of the child element's keys
        if parent.Elem[i].Key != nil {
            for key, value := range parent.Elem[i].Key {
                if child.Elem[i].Key == nil || child.Elem[i].Key[key] != value {
                    return false
                }
            }
        }
    }

    return true
}


// ParsePathFromString parses a string representation of a path and creates a gnmi.Path object.
// func ParsePathFromString(pathStr string) (*gnmi.Path, error) {
//     // Split the input string by "elem:" to extract individual path elements
//     pathElemStrs := strings.Split(pathStr, "elem:")

//     // Initialize a gnmi.Path object
//     path := &gnmi.Path{
//         Elem: make([]*gnmi.PathElem, 0),
//     }

//     // Iterate over path elements and parse them
//     for _, pathElemStr := range pathElemStrs {
//         // Skip empty strings
// 		if pathElemStr == "" {
// 			continue
// 		}
//         // Trim spaces and curly braces
// //		pathWord := extractWords(pathElemStr)

//         // Initialize a gnmi.PathElem object
//         pathElem := &gnmi.PathElem{
//             Key: make(map[string]string),
//         }

//         // Split the path element string by whitespace
//         elemParts := strings.Fields(pathElemStr)

//         for _, elemPart := range elemParts {
//             // Split each part by ":" to extract the field name and value
//             kv := strings.SplitN(elemPart, ":", 2)
//             if len(kv) != 2 {
//                 return nil, fmt.Errorf("Invalid element format: %s", elemPart)
//             }

//             key := strings.TrimSpace(kv[0])
//             value := strings.TrimSpace(kv[1])

//             key = extractWords(key)
//             value = extractWords(value)
//             switch key {
//             case "name":
//                 pathElem.Name = value
//             case "key":
//                 // Parse key field and value pairs
//                 keyValues := strings.Split(value, ",")
//                 for _, keyValueStr := range keyValues {
//                     keyValue := strings.Split(keyValueStr, ":")
//                     if len(keyValue) != 2 {
//                         return nil, fmt.Errorf("Invalid key format: %s", keyValueStr)
//                     }
//                     pathElem.Key[strings.TrimSpace(keyValue[0])] = strings.TrimSpace(keyValue[1])
//                 }
//             }
//         }

//         // Add the parsed path element to the gnmi.Path
//         path.Elem = append(path.Elem, pathElem)
//     }

//     return path, nil
// }



func ParsePathFromString(pathStr string) (*gnmi.Path, error) {
	// Use regular expressions to match and extract elements
	elemRegex := regexp.MustCompile(`elem:{([^}]+)}`)
	elemMatches := elemRegex.FindAllStringSubmatch(pathStr, -1)

	// Initialize a gnmi.Path object
	path := &gnmi.Path{
		Elem: make([]*gnmi.PathElem, 0),
	}

	for _, elemMatch := range elemMatches {
		elemStr := elemMatch[1]

		// Initialize a gnmi.PathElem object
		pathElem := &gnmi.PathElem{
			Key: make(map[string]string),
		}

		// Split the element string by ","
		elemParts := strings.Split(elemStr, ",")

		for _, elemPart := range elemParts {
			// Split each part by ":" to extract the field name and value
			kv := strings.SplitN(elemPart, ":", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("Invalid element format: %s", elemPart)
			}

			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])

			// Remove escaped double quotes from value
			value = strings.Replace(value, `\"`, `"`, -1)

			if key == "name" {
				pathElem.Name = value
			} else {
				// Handle key field and value pairs
				pathElem.Key[key] = value
			}
		}

		// Add the parsed path element to the gnmi.Path
		path.Elem = append(path.Elem, pathElem)
	}

	return path, nil
}

/// MergePaths merges a child path into a parent path based on the described rules.
func MergePaths(parentPath, childPath *gnmi.Path) *gnmi.Path {
    mergedElem := make([]*gnmi.PathElem, len(childPath.Elem))
    
    for i, childElem := range childPath.Elem {
        var mergedElemName string
        
		if i >= len(parentPath.Elem) {
			mergedElem[i] = childElem
			continue
		}
        // If the names match, use the child's name
        if i < len(parentPath.Elem) && childElem.Name == parentPath.Elem[i].Name {
            mergedElemName = childElem.Name
        } else {
            mergedElemName = parentPath.Elem[i].Name
        }
        
        // Copy the child element
        mergedElem[i] = &gnmi.PathElem{
            Name: mergedElemName,
        }
        
        // If the child element has keys and the parent does not, use the child's keys
        if len(childElem.Key) > 0 && (i >= len(parentPath.Elem) || len(parentPath.Elem[i].Key) == 0) {
            mergedElem[i].Key = make(map[string]string)
            for key, value := range childElem.Key {
                mergedElem[i].Key[key] = value
            }
        }
    }
    
    return &gnmi.Path{
        Elem: mergedElem,
    }
}

func main() {
     // Create a gNMI Path object
     path := &gnmi.Path{
        Elem: []*gnmi.PathElem{
            {Name: "interfaces"},
            {Name: "interface", Key: map[string]string{"name": "eth0"}},
        },
    }

    // Get the string representation of the Path object
    pathStr := path.String()

    // Print the string representation
    fmt.Println(pathStr)

    childPathString := "elem:{name:\"interfaces\"}  elem:{name:\"interface\"  key:{key:\"name\"  value:\"admin4\"}}"
	parentPathString := "elem:{name:\"openconfig-interfaces:interfaces\"}  elem:{name:\"interface\"}"
	childPath, _ := ParsePathFromString(childPathString)
	parentPath, _ := ParsePathFromString(parentPathString)

	if pathContains(parentPath, childPath) {
        fmt.Println("Parent path 2 contains the child path 2.")
    } else {
        fmt.Println("Parent path 2 does not contain the child path 2.")
    }

    mergePath := MergePaths(parentPath, childPath)
    fmt.Println(mergePath)
}