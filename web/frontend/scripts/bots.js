let rotateableImg = document.querySelector('.rotatable-image');


let isDragging = false;
let previousMousePosition = { x:0, y:0};
let rotation = 0;


function getRelativeCoordinates(event) {
    const rect = rotateableImg.getBoundingClientRect();
    
    // Calculate the center of the image
    const centerX = rect.left + rect.width / 2;
    const centerY = rect.top + rect.height / 2;

    // Calculate the mouse position relative to the image's center
    const relativeX = event.clientX - centerX;
    const relativeY = event.clientY - centerY;

    return { relativeX, relativeY };
}

function getDistanceFromCenter(relativeX, relativeY) {
    return Math.sqrt(relativeX * relativeX + relativeY * relativeY);
}



rotateableImg.addEventListener('mousedown', (event) => {
    isDragging = true;
    previousMousePosition = getRelativeCoordinates(event);

});

rotateableImg.addEventListener('mouseleave', (event) => {
    isDragging = false;
});


document.addEventListener('mousemove', (event) => {

    if (!isDragging) return;

    // Get the current mouse position relative to the center of the image
    const currentMousePosition = getRelativeCoordinates(event);

    // Calculate the angle change
    const angleStart = Math.atan2(previousMousePosition.relativeY, previousMousePosition.relativeX);
    const angleEnd = Math.atan2(currentMousePosition.relativeY, currentMousePosition.relativeX);
    
    // Determine the rotation direction
    let angleDelta = angleEnd - angleStart;

    // Normalize the angle to get the direction
    if (angleDelta > Math.PI) {
        angleDelta -= 2 * Math.PI;
    } else if (angleDelta < -Math.PI) {
        angleDelta += 2 * Math.PI;
    }
    const distanceFromCenter = getDistanceFromCenter(currentMousePosition.relativeX, currentMousePosition.relativeY);
    // Adjust the rotation
    rotation += angleDelta * distanceFromCenter * 0.3;
    // Apply the rotation
    rotateableImg.style.transform = `rotate(${rotation}deg)`;

    // Update the previous mouse position
    previousMousePosition = currentMousePosition;
});

rotateableImg.addEventListener('mouseup', (event) => {

    isDragging = false;
});