function toggleTitleAnimation() {

    // console.log('Tittle effect toggle');
    const title = document.getElementById('hof-title');

    if (title.style.animationPlayState === 'paused') {
        title.style.animationPlayState = 'running'; // Resume the animation
    } else {
        title.style.animationPlayState = 'paused';  // Pause the animation (freeze)
    }
}


function createMatrixEffect() {
    let matrixActive = false; // Internal state for the matrix effect
    let intervalID = null;   // Holds the interval ID

    const canvas = document.getElementById('matrixCanvas');
    const context = canvas.getContext('2d');

    const mainContent = document.getElementById("main-content");

    // Set canvas size based on parent element
    canvas.width = mainContent.offsetWidth;
    canvas.height = mainContent.offsetHeight;

    const katakana = 'アァカサタナハマヤャラワガザダバパイィキシチニヒミリヰギジヂビピウゥクスツヌフムユュルグズブヅプエェケセテネヘメレヱゲゼデベペオォコソトノホモヨョロヲゴゾドボポヴッン';
    const latin = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
    const nums = '0123456789';
    const chars = '$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$';
    const alphabet = katakana + latin + nums + chars;

    const fontSize = 16;
    const columns = canvas.width / fontSize;

    const rainDrops = [];
    for (let x = 0; x < columns; x++) {
        rainDrops[x] = 1;
    }

    const draw = () => {
        if (!matrixActive) {
            return; // Stop drawing if matrixActive is false
        }

        context.fillStyle = 'rgba(245, 242, 235, 0.085)';
        context.fillRect(0, 0, canvas.width, canvas.height);

        // Letters style
        context.fillStyle = 'rgb(255, 209, 109)';
        context.font = fontSize + 'px monospace';

        for (let i = 0; i < rainDrops.length; i++) {
            const text = alphabet.charAt(Math.floor(Math.random() * alphabet.length));
            context.fillText(text, i * fontSize, rainDrops[i] * fontSize);

            if (rainDrops[i] * fontSize > canvas.height && Math.random() > 0.975) {
                rainDrops[i] = 0;
            }
            rainDrops[i]++;
        }
    };

    // Function to start the matrix effect
    function startMatrixEffect() {
        if (!matrixActive) {
            matrixActive = true;
            intervalID = setInterval(draw, 60);
            // console.log("Matrix effect started");
        }
    }

    // Function to stop the matrix effect
    function stopMatrixEffect() {
        if (matrixActive) {
            clearInterval(intervalID);
            matrixActive = false;
            // console.log("Matrix effect stopped");
        }
    }

    // Function to toggle the matrix effect
    function toggleMatrixEffect(event) {
        if (event) {
            event.stopPropagation(); // Stop event propagation
        }
        if (matrixActive) {
            stopMatrixEffect();
        } else {
            startMatrixEffect();
        }
    }

    // Automatically start the matrix effect upon initialization
    if (!intervalID) {
        // console.log("Starting matrix effect...");
        startMatrixEffect();
    }
    // Return an interface to control the matrix effect externally
    return {
        start: startMatrixEffect,
        stop: stopMatrixEffect,
        toggle: toggleMatrixEffect
    };
}


function setHofTitle(guild_name, month, year) {

    const matrix_switch = document.getElementById('matrix-switch');
    const title_effect_switch = document.getElementById('title-effect-switch');

    matrix_switch.setAttribute('checked', 'checked');
    title_effect_switch.setAttribute('checked', 'checked');

    if (!month || !year) {
        let currentDate = new Date();
        month = currentDate.getMonth() + 1;
        year = currentDate.getFullYear();
    }

    let hofTitle = document.getElementById('hof-title');
    let currentDate = new Date();

    currentDate.setMonth(currentDate.getMonth() - 1);

    hofTitle.innerText = `${guild_name} HALL OF FAME - ${month}/${year}`;
}

