$(function() {
  function setMessage(message) {
    $("#message").text(message);
  }

  var ws = new WebSocket("ws://" + location.host + "/connect");
  ws.onopen = function() {
    console.log("OPEN");
    $("body").addClass("connected");
    setMessage("connecting...");
  };
  ws.onclose = function() {
    console.log("CLOSE");
    $("body").removeClass("connected");
    setMessage("closed.");
  };
  ws.onmessage = function(event) {
    console.log(event);

    var xml = $.parseXML(event.data).firstChild;
    var docType = xml.nodeName;
    var handler = handlers[docType];
    if (handler) {
      new handler(xml);
    } else {
      console.log("unknown element type");
      setMessage("unknown message");
    }
  };

  var handlers = {};
  handlers.DisplayLayouts = function(xml) {
    setMessage("layout");
    window.d = xml;

    var layouts = {};
    eachChildElement(xml, function(layoutNode) {
      var layout = new Layout(layoutNode);
      layouts[layout.name] = layout;
    });
    layouts[xml.getAttribute("selected")].arrange($("#stage-container"));
  };
  handlers.StageDisplayData = function(xml) {
    setMessage("slide");
    window.s = xml;

    var fields = xml.firstElementChild;
    eachChildElement(fields, function(field) {
      $("." + classForField(field)).text($(field).text());
    });

    console.log("New slide");
    console.log(xml);
  };

  function eachChildElement(xml, fn) {
    for (var i = 0; i < xml.childNodes.length; i++) {
      fn(xml.childNodes[i]);
    }
  }

  function Layout(xml) {
    this.name = xml.getAttribute("identifier");

    this.arrange = function($e) {
      $e.empty().append(build());
    };

    //var totalWidth = xml.getAttribute("width");
    //var totalHeight = xml.getAttribute("height");

    function build() {
      var main = buildBlock(xml);
      main.addClass("display-frame");
      main.css("position", "relative");
      eachChildElement(xml, function(frame) {
        main.append(buildFrame(frame));
      });
      return main;
    };

    function buildFrame(node) {
      var block = buildBlock(node);
      var identifier = node.getAttribute("identifier");
      block.css("position", "absolute");
      block.css("top", coord(node, "yAxis"));
      block.css("left", coord(node, "xAxis"));
      block.text(identifier);
      var label = $("<div></div>").text(identifier).css("position", "absolute").css("top", 0).css("left", 0).css("background-color", "black").css("color", "white").css("font-weight", "bold");
      var content = $("<div></div>").css("position", "absolute").css("top", 0).css("bottom", 0).css("left", 0).css("right", 0);
      content.addClass(classForField(node));
      content.css("font-size", node.getAttribute("fontSize"));
      block.append(label).append(content);
      return block;
    }

    function buildBlock(node) {
      var e = $("<div></div>");
      e.width(coord(node, "width"));
      e.height(coord(node, "height"));
      e.css("border", "solid 1px black");
      return e;
    }

    function coord(node, attrName) {
      return parseInt(node.getAttribute(attrName));
    }

    return this;
  }

  function classForField(node) {
    return "js-layout-" + node.getAttribute("identifier");
  }
});
// <DisplayLayouts selected="Default">
//   <DisplayLayout showBorder="1E0" width="1680" identifier="Default" height="1050">
//     <Frame height="105.000000" width="336.000000" xAxis="84.000000" isVisible="YES" identifier="Clock" yAxis="0.000000"/>
//     <Frame height="105.000000" width="336.000000" xAxis="1260.000000" isVisible="YES" identifier="ElapsedTime" yAxis="0.000000"/>
//     <Frame height="105.000000" width="336.000000" xAxis="84.000000" isVisible="YES" identifier="Timer1" yAxis="787.500000"/>
//     <Frame height="105.000000" width="336.000000" xAxis="672.000000" isVisible="YES" identifier="Timer2" yAxis="787.500000"/>
//     <Frame height="105.000000" width="336.000000" xAxis="1260.000000" isVisible="YES" identifier="Timer3" yAxis="787.500000"/>
//     <Frame height="105.000000" width="336.000000" xAxis="672.000000" isVisible="YES" identifier="VideoCounter" yAxis="0.000000"/>
//     <Frame height="420.000000" width="336.000000" xAxis="1302.000000" isVisible="YES" identifier="ChordChart" yAxis="236.250000"/>
//     <Frame height="525.000000" width="672.000000" xAxis="42.000000" isVisible="YES" identifier="CurrentSlide" yAxis="131.250000" fontSize="60"/>
//     <Frame height="420.000000" width="504.000031" xAxis="756.000000" isVisible="YES" identifier="NextSlide" yAxis="183.750000" fontSize="60"/>
//     <Frame height="105.000000" width="672.000000" xAxis="42.000000" isVisible="YES" identifier="CurrentSlideNotes" yAxis="656.250000" fontSize="60"/>
//     <Frame height="105.000000" width="504.000031" xAxis="756.000000" isVisible="YES" identifier="NextSlideNotes" yAxis="603.750000" fontSize="60"/>
//     <Frame height="105.000000" width="1512.000000" xAxis="84.000000" isVisible="YES" identifier="Message" yAxis="918.750000" fontSize="60" flashColor="0.000000 1.000000 0.000000"/>
//   </DisplayLayout>
// </DisplayLayouts>


// <?xml version="1.0" encoding="UTF-8" standalone="no"?>
// <StageDisplayData>
//   <Fields>
//     <Field type="clock" clockFormat="1" label="Clock" identifier="Clock" alpha="1E0" red="1E0" green="1E0" blue="1E0">7:49:48 PM</Field>
//     <Field running="0" type="elapsed" label="Time Elapsed" identifier="ElapsedTime" alpha="1E0" red="1E0" green="1E0" blue="1E0">--:--:--</Field>
//     <Field running="0" type="countdown" overrun="0" label="Countdown 1" identifier="Timer1" alpha="1E0" red="1E0" green="1E0" blue="1E0">--:--:--</Field>
//     <Field running="0" type="countdown" overrun="0" label="Countdown 2" identifier="Timer2" alpha="1E0" red="1E0" green="1E0" blue="1E0">--:--:--</Field>
//     <Field running="0" type="countdown" overrun="0" label="Countdown 3" identifier="Timer3" alpha="1E0" red="1E0" green="1E0" blue="1E0">--:--:--</Field>
//     <Field type="slide" label="Current Slide" identifier="CurrentSlide" alpha="1E0" red="1E0" green="1E0" blue="1E0">You're the reason I sing
// The reason I sing
// Yes my heart will sing
// How I love You</Field>
//     <Field type="slide" label="Next Slide" identifier="NextSlide" alpha="1E0" red="1E0" green="1E0" blue="1E0">And forever I'll sing
// Forever I'll sing
// Yes my heart will sing
// How I love You</Field>
//     <Field type="slide" label="Current Slide Notes" identifier="CurrentSlideNotes" alpha="1E0" red="1E0" green="1E0" blue="1E0"/>
//     <Field type="slide" label="Next Slide Notes" identifier="NextSlideNotes" alpha="1E0" red="1E0" green="1E0" blue="1E0"/>
//     <Field type="message" label="Message" identifier="Message" alpha="1E0" red="1E0" green="1E0" blue="1E0"/>
//     <Field running="0" type="countdown" label="Video Countdown" identifier="VideoCounter" alpha="1E0" red="1E0" green="1E0" blue="1E0">--:--:--</Field>
//     <Field type="chordChart" label="Chord Chart" identifier="ChordChart"/>
//   </Fields>
// </StageDisplayData>
