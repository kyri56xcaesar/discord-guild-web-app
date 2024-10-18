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

// Create and initialize the matrix effect
matrixEffect = createMatrixEffect();

// Now you can control the matrix effect via the returned object
// Example usage for toggling the matrix effect
document.getElementById('matrix-switch').addEventListener('click', matrixEffect.toggle);


