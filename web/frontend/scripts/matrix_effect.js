matrixActive = true;

const canvas = document.getElementById('matrixCanvas');
const context = canvas.getContext('2d');

canvas.width = window.innerWidth;
canvas.height = window.innerHeight;



const katakana = 'アァカサタナハマヤャラワガザダバパイィキシチニヒミリヰギジヂビピウゥクスツヌフムユュルグズブヅプエェケセテネヘメレヱゲゼデベペオォコソトノホモヨョロヲゴゾドボポヴッン';
const latin = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
const nums = '0123456789';
const chars = '$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$'

const alphabet = katakana + latin + nums + chars;

const fontSize = 16;
const columns = canvas.width/fontSize;

const rainDrops = [];

for( let x = 0; x < columns; x++ ) {
	rainDrops[x] = 1;
}

const draw = () => {
	context.fillStyle = 'rgba(245, 242, 235, 0.085)';
	context.fillRect(0, 0, canvas.width, canvas.height);
	
	context.fillStyle = 'rgb(255, 209, 109)';
	context.font = fontSize + 'px monospace';

	for(let i = 0; i < rainDrops.length; i++)
	{
		const text = alphabet.charAt(Math.floor(Math.random() * alphabet.length));
		context.fillText(text, i*fontSize, rainDrops[i]*fontSize);
		
		if(rainDrops[i]*fontSize > canvas.height && Math.random() > 0.975){
			rainDrops[i] = 0;
        }
		rainDrops[i]++;
	}
};

if (matrixActive) {
    setInterval(draw, 60);
}



function toggleMatrixEffect() {
    console.log(matrixActive)
    if (matrixActive) {
        matrixActive = false
    } else {
        matrixActive = true;
    }
}