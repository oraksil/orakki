window.addEventListener("resize", resizeCanvas, false);

// const target =
// document.getElementById("video") || document.getElementById("canvas");
const target = document.getElementById("gamescreen");
function resizeCanvas() {
  if (!target) return;
  var width = window.innerWidth || document.body.clientWidth;
  var height = window.innerHeight || document.body.clientHeight;

  if (height > width) {
    target.style["width"] = "80%";
    target.style["height"] = "";
  } else {
    target.style["height"] = height + "px"
    target.style["width"] = "";
  }
}

resizeCanvas();
