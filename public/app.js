function setMessage(message) {
  document.getElementById('message').textContent = message;
}

var ws = new WebSocket("ws://" + location.host + "/connect");
ws.onopen = function() {
  console.log("OPEN");
  document.body.style.backgroundColor = '#cfc';
  setMessage("connecting...");
};
ws.onclose = function() {
  console.log("CLOSE");
  document.body.style.backgroundColor = null;
};
ws.onmessage = function(event) {
  console.log(event);
  setMessage(event.data);
};

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
