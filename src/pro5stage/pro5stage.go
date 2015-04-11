package pro5stage

import (
  "bytes"
  "container/list"
  "encoding/xml"
  "fmt"
  "io"
  "log"
  "net"
  "time"
)

// Arguments for connecting to ProPresenter.
type ConnectInfo struct {
  Host string
  Port int
  Password string
}

// A connection to ProPreseter.
type Conn struct {
  Info ConnectInfo
  Connection net.Conn
  Listeners list.List
  DisplayLayouts string
  LastSlide string
}

func ConnectToPro5(info ConnectInfo) (*Conn, error) {
  var result = new(Conn)
  var err error
  result.Info = info
  result.Connection, err = net.Dial("tcp", fmt.Sprintf("%s:%d", info.Host, info.Port))
  if err != nil {
    return nil, err
  }
  go result.Run()
  return result, nil
}

func (c *Conn) Run() {
  var xmlWriter = xml.NewEncoder(c.Connection)
  var loginElement = xml.StartElement{}
  loginElement.Name.Local = "StageDisplayLogin"
  xmlWriter.EncodeElement(c.Info.Password, loginElement)
  fmt.Fprintf(c.Connection, "\r\n")

  c.ReadEverything()
}

func (c *Conn) ReadEverything() {
  var xmlReader = xml.NewDecoder(c.Connection)

  // http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/
  for {
    var token, err = xmlReader.Token()
    if err != nil {
      log.Fatal("xmlReader.Token(): ", err)
    }
    switch se := token.(type) {
    case xml.StartElement:
      switch se.Name.Local {
      case "DisplayLayouts":
        c.DisplayLayouts, err = readXmlString(xmlReader, &se)
        if err != nil {
          log.Fatal(err)
        }
        sendToListeners(c, c.DisplayLayouts)
      case "StageDisplayData":
        c.LastSlide, err = readXmlString(xmlReader, &se)
        if err != nil {
          log.Fatal(err)
        }
        sendToListeners(c, c.LastSlide)
      }
    }
  }
}

func readXmlString(xmlReader *xml.Decoder, startElement *xml.StartElement) (string, error) {
  var err error

  var buffer bytes.Buffer
  var xmlWriter = xml.NewEncoder(&buffer)

  err = xmlWriter.EncodeToken(*startElement)
  if err != nil {
    return "", err
  }
  for depth := 1; depth > 0; {
    var token, err = xmlReader.Token()
    if err != nil {
      return "", err
    }
    switch token.(type) {
    case xml.StartElement:
      depth = depth + 1
    case xml.EndElement:
      depth = depth - 1
    }
    err = xmlWriter.EncodeToken(token)
    if err != nil {
      return "", err
    }
  }
  xmlWriter.Flush()
  return buffer.String(), nil
}

func (c *Conn) AddListener(listener io.Writer) {
  c.Listeners.PushBack(listener)
  sendToListener(c, listener, c.DisplayLayouts)
  sendToListener(c, listener, c.LastSlide)
  for {
    time.Sleep(1 * time.Minute)
  }
}

func sendToListeners(c *Conn, payload string) (err error) {
  for e := c.Listeners.Front(); e != nil; e = e.Next() {
    listener, ok := e.Value.(io.Writer)
    if ok {
      err = sendToListener(c, listener, payload)
      if err != nil {
        return
      }
    }
  }
  return
}

func sendToListener(c *Conn, listener io.Writer, payload string) error {
  if len(payload) > 0 {
    _, err := fmt.Fprintf(listener, "%s", payload)
    if err != nil {
      return err
    }
  }
  return nil
}
