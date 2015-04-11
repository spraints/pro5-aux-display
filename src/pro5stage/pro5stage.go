package pro5stage

import (
  "bytes"
  "encoding/xml"
  "fmt"
  "log"
  "net"
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
}

type Client interface {
  SendMessage(name string, payload string)
}

func Run(connectInfo ConnectInfo, client Client) {
  pro5, err := ConnectToPro5(connectInfo)
  if err != nil {
    log.Fatal("ConnectToPro5: ", err)
  }
  pro5.ReadEverything(client)
}

func ConnectToPro5(info ConnectInfo) (*Conn, error) {
  var result = new(Conn)
  var err error
  result.Info = info
  result.Connection, err = net.Dial("tcp", fmt.Sprintf("%s:%d", info.Host, info.Port))
  if err != nil {
    return nil, err
  }
  connect(result)
  return result, nil
}

func connect(c *Conn) {
  var xmlWriter = xml.NewEncoder(c.Connection)
  var loginElement = xml.StartElement{}
  loginElement.Name.Local = "StageDisplayLogin"
  xmlWriter.EncodeElement(c.Info.Password, loginElement)
  fmt.Fprintf(c.Connection, "\r\n")
}

func (c *Conn) ReadEverything(client Client) {
  var xmlReader = xml.NewDecoder(c.Connection)

  // http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/
  for {
    var token, err = xmlReader.Token()
    if err != nil {
      log.Fatal("xmlReader.Token(): ", err)
    }
    switch se := token.(type) {
    case xml.StartElement:
      xml, err := readXmlString(xmlReader, &se)
      if err != nil {
        log.Fatal(err)
      }
      client.SendMessage(se.Name.Local, xml)
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
