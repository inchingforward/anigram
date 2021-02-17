var Editor = Editor || (function() {
    var PIXEL_SIZE = 16;
    var SIDE = PIXEL_SIZE * PIXEL_SIZE;

    var DARK = "#000000";
    var MEDIUM = "#999999";
    var BACKGROUND = "#d2d2d2";
    var LIGHT = "#ffffff";
    
    var canvas, ctx;
    var backgroundBox, lightBox, mediumBox, darkBox, currentBox;
    var frameText;
    var color = MEDIUM;
    var mouseDown = false;
    
    var animation, currFrameIndex, animation_uuid;
    var emptyFrame = "0".repeat(SIDE);

    function getAnimation() {
        return animationToString(animation);
    }

    function animationToString(animationArray) {
        var animationStr = "";
        
        for (var i=0; i < animationArray.length; i++) {
            animationStr += animationArray[i];
        }

        return animationStr;
    }

    function stringToAnimation(str) {
        if (!str) {
            return [emptyFrame];
        }

        var animArray = [];

        for (var i=0; i < str.length; i += SIDE) {
            animArray.push(str.substring(i, i + SIDE));
        }

        return animArray;
    }

    function updateFrameText() {
        frameText.innerHTML = "Frame " + (currFrameIndex + 1) + " of " + animation.length;
    }

    function prevFrame() {
        var nextFrameIndex = currFrameIndex - 1;
        if (nextFrameIndex < 0) {
            return;
        }

        currFrameIndex = nextFrameIndex;

        loadFrame(animation[currFrameIndex]);
        updateFrameText();
    }

    function nextFrame() {
        if (currFrameIndex >= animation.length - 1) {
            animation.push(emptyFrame);
        }

        currFrameIndex++;
        
        loadFrame(animation[currFrameIndex]);
        updateFrameText();
    }

    function clear() {
        ctx.beginPath();
        ctx.fillStyle = BACKGROUND;
        ctx.rect(0, 0, SIDE, SIDE);
        ctx.fill();

        animation[currFrameIndex] = emptyFrame;
    }

    function play() {
        currFrameIndex = 0;

        var intervalId = setInterval(function() {
            loadFrame(animation[currFrameIndex]);
            updateFrameText();

            if (currFrameIndex == (animation.length - 1)) {
                clearInterval(intervalId);
                return;
            }
            
            currFrameIndex += 1;
        }, 325);
    }

    function save() {
        var animationStr = animationToString(animation);
        
        // FIXME
        axios.patch("/api/animations/" + animation_uuid, {
                animation: animationStr
            }).then(function(response) {
                console.log("success", response);
            }).catch(function(error) {
                console.log("error", error);
            });
    }

    function setColor(nextColor) {
        color = nextColor;
        currentBox.style.backgroundColor = color;
    }

    function canvasMouseDown(e) {
        mouseDown = true;
        canvasMouseOver(e);
    }

    function canvasMouseMove(e) {
        canvasMouseOver(e);
    }

    function canvasMouseOver(e) {
        if (!mouseDown) {
            return;
        }

        var boundingRect = canvas.getBoundingClientRect();
        var x = e.clientX - boundingRect.left;
        var y = e.clientY - boundingRect.top;

        var row = Math.floor(y / PIXEL_SIZE);
        var col = Math.floor(x / PIXEL_SIZE);

        fillSquareAt(row, col, color, true);
    }

    function charToColor(char) {
        if (char === "0") {
            return BACKGROUND;
        } else if (char === "1") {
            return LIGHT;
        } else if (char === "2") {
            return MEDIUM;
        } else if (char === "3") {
            return DARK;
        } else {
            return BACKGROUND;
        }
    }

    function colorToChar(color) {
        if (color === BACKGROUND) {
            return "0";
        } else if (color === LIGHT) {
            return "1";
        } else if (color === MEDIUM) {
            return "2";
        } else if (color === DARK) {
            return "3";
        } else {
            return "0";
        }
    }

    function indexToCoord(index) {
      return [index % PIXEL_SIZE, Math.floor(index / PIXEL_SIZE)];
    }

    function coordToIndex(row, col) {
        return row * PIXEL_SIZE + col;
    }
    
    function loadFrame(frame) {
        // FIXME: bump the frame number display.
        for (var i=0; i < frame.length; i++) {
            var ch = frame[i];
            var color = charToColor(ch);
            var coord = indexToCoord(i);
            
            fillSquareAt(coord[1], coord[0], color, false);
        }
    }

    function updateFrameSquare(row, col, color) {
        var frame = animation[currFrameIndex];
        var frameIndex = coordToIndex(row, col);
        var color = colorToChar(color);

        frame = animation[currFrameIndex];
        animation[currFrameIndex] = frame.slice(0, frameIndex) + color + frame.slice(frameIndex + 1); 
    }

    function fillSquareAt(row, col, color, updateAnimation) {
        ctx.beginPath();
        ctx.fillStyle = color;
        ctx.rect(col * PIXEL_SIZE, row * PIXEL_SIZE, PIXEL_SIZE, PIXEL_SIZE);
        ctx.fill();

        if (updateAnimation) {
            updateFrameSquare(row, col, color);
        }
    }

    function keyUp(e) {
        if (e.code === "ArrowLeft") {
            prevFrame();
        } else if (e.code === "ArrowRight") {
            nextFrame();
        } else if (e.code === "KeyC") {
            clear();
        } else if (e.code === "KeyP") {
            play();
        }
    }

    function mouseUp(e) {
        mouseDown = false;
    }

    function dumpAnimation() {
        console.log(animation);
    }
    
    function init(uuid) {
        canvas = canvas || document.getElementById("canvas-grid");
        ctx = ctx || canvas.getContext("2d");

        canvas.width = SIDE;
        canvas.height = SIDE;

        ctx.fillStyle = BACKGROUND;
        ctx.beginPath();
        ctx.rect(0, 0, SIDE, SIDE);
        ctx.fill();
        
        backgroundBox = document.getElementById("background-box");
        backgroundBox.style.backgroundColor = BACKGROUND;
        backgroundBox.onclick = function() { setColor(BACKGROUND); };

        lightBox = document.getElementById("light-box");
        lightBox.style.backgroundColor = LIGHT;
        lightBox.onclick = function() { setColor(LIGHT); };

        mediumBox = document.getElementById("medium-box");
        mediumBox.style.backgroundColor = MEDIUM;
        mediumBox.onclick = function() { setColor(MEDIUM); };

        darkBox = document.getElementById("dark-box");
        darkBox.style.backgroundColor = DARK;
        darkBox.onclick = function() { setColor(DARK); };

        currentBox = document.getElementById("current-box");
        
        setColor(MEDIUM);

        document.getElementById("prevFrame").onclick = prevFrame;
        document.getElementById("nextFrame").onclick = nextFrame;
        document.getElementById("clear").onclick = clear;
        document.getElementById("play").onclick = play;
        document.getElementById("save").onclick = save;

        canvas.addEventListener("mousemove", canvasMouseMove);
        canvas.addEventListener("mousedown", canvasMouseDown);
        window.addEventListener("mouseup", mouseUp);
        document.addEventListener("keyup", keyUp);

        frameText = document.getElementById("frameText");

        if (!uuid) {
            animation = [];
            animation.push(emptyFrame);
            currFrameIndex = 0;
            loadFrame(animation[0]);
        } else {
            animation_uuid = uuid;

            /*axios.get("/api/animations/" + animation_uuid
            ).then(function(response) {
                data = response.data;

                if (!data.animation) {
                    // FIXME:  What should we do here?
                    alert("Bad or corrupt animation!");
                    return;
                }*/
                var data = {
                    title: "Testing",
                    animation: "000000000000000000000000000000000010010000000000002002000000000000200200000000000020020100000000002002020000000000200200000000000021120100000000002222020000000000200202000000000020020200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000222220000000000222222200000002222000022000000020000002220000000000000222000000000000222200000000022222000000000222222000000000220000000000000022000000000000022200002222000002222222222000000022222222000000000000000000000000000000000000000000000000000000000022200000000000022222200000000222000220000000022200002200000002200000220000000000002222000000000002222000000000000000200000000000000022000000022200022200000002220002220000000222222222000000000022222000000000000000000000000000000000000000000000000000000000000022000000000000022200000000000022200000000000022220000000000022202000000000002200200000000002220020000000000222222222200000002222222220000000000020000000000000002000000000000002220000000000000222000000000000022200000000000000000000000000000000000000000022222222200000002222222220000002220000022000000222000000000000022222222200000000222222222000000022220002220000000000000222000022000000022200002220000002220000222000000220000022222002222000000222222222000000000000000000000000000000000000000000000000000000000002222200000000022222222000000022222220000000022000000000000002200000000000000220022222000000022022222220000002222000022200000222000000220000002200000222000000222002222200000002222222200000000222222200000000002222000000000000000000000000000000000000000000000000000000000022222222220000022222222222000002222200022000000022000002200000000000002200000000000002220000000000000220000000000000022000000000000022000000000000022200000000000002220000000000002222000000000000022200000000000000000000000000000000000000000000222200000000000222222000000000022222220000000002200022000000000220002200000000022222200000000000222220000000000222202200000000220000222000000022000022200000002200002220000000222222220000000022222220000000000000000000000000000000000000000000000000000000000022222200000000002222220000000002222222200000002220000220000000222000022000000022200002200000000220000220000000022200222000000002222222200000000022222020000002200000022000000222000022000000002222222200000000022222000000000000000000000000000000000000000000000000000000000002000022220000002200022222200002220002222220000022002200002200002200200000220000220020000022000022002000002200002200200000220000220022000220000022002222222000022220222222000002222002222000000222000000000000000000000000000"
                }

                title = data.title;
                animation = stringToAnimation(data.animation);
                currFrameIndex = 0;
                
                loadFrame(animation[0]);

                updateFrameText();
            /*}).catch(function(error) {
                console.log("error", error);
                alert("Could not load animation");
            });*/
        }
    }

    return {
        init: init,
        dumpAnimation: dumpAnimation,
        getAnimation: getAnimation
    }
})();