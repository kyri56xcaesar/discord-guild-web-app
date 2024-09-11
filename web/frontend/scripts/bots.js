document.addEventListener("DOMContentLoaded", function() {
    const rotatableImage = document.querySelector('.rotatable-image');
    let isDragging = false;
    let previousAngle = 0;
    let startX, startY, currentRotation = 0, rotationVelocity = 0, deceleration = 0.99;
    let lastFrameTime = 0, lastMouseMoveTime = 0;
    const accelerationFactor = 0.99; // Adjust this factor to increase the acceleration effect

    function getAngle(x, y) {
        return Math.atan2(y, x) * (180 / Math.PI);
    }

    rotatableImage.addEventListener('mousedown', function(event) {
        isDragging = true;
        const rect = rotatableImage.getBoundingClientRect();
        startX = event.clientX - rect.left - rect.width / 2;
        startY = event.clientY - rect.top - rect.height / 2;
        previousAngle = getAngle(startX, startY) - currentRotation;  // Adjust for current rotation
        rotatableImage.style.cursor = 'grabbing';
        rotationVelocity = 0;  // Reset velocity on drag start
        lastMouseMoveTime = performance.now();  // Track time to calculate velocity
    });

    document.addEventListener('mousemove', function(event) {
        if (isDragging) {
            const rect = rotatableImage.getBoundingClientRect();
            const currentX = event.clientX - rect.left - rect.width / 2;
            const currentY = event.clientY - rect.top - rect.height / 2;
            const currentAngle = getAngle(currentX, currentY);
            const currentTime = performance.now();
            const timeDiff = currentTime - lastMouseMoveTime;  // Time between frames

            // Accumulate velocity with acceleration
            const angleDifference = currentAngle - previousAngle;
            const velocityBoost = angleDifference * (1000 / timeDiff) * accelerationFactor;
            rotationVelocity += velocityBoost;  // Add to velocity based on speed of dragging

            // Update rotation and previous values
            currentRotation = currentAngle - previousAngle;
            rotatableImage.style.transform = `rotate(${currentRotation}deg)`;

            lastMouseMoveTime = currentTime;
        }
    });

    document.addEventListener('mouseup', function(event) {
        if (isDragging) {
            isDragging = false;
            rotatableImage.style.cursor = 'grab';

            // Calculate rotation speed
            const rect = rotatableImage.getBoundingClientRect();
            const endX = event.clientX - rect.left - rect.width / 2;
            const endY = event.clientY - rect.top - rect.height / 2;
            const endAngle = getAngle(endX, endY);
            rotationVelocity += (endAngle - previousAngle) * accelerationFactor;  // Final velocity boost

            // Start the spinning animation based on velocity
            lastFrameTime = performance.now();
            requestAnimationFrame(animateSpin);
        }
    });

    function animateSpin(timestamp) {
        const deltaTime = timestamp - lastFrameTime;
        lastFrameTime = timestamp;

        // Apply deceleration to slow the spin over time
        rotationVelocity *= deceleration;
        currentRotation += rotationVelocity * (deltaTime / 16);  // Scale with deltaTime for consistent speed

        // Apply rotation to the image
        rotatableImage.style.transform = `rotate(${currentRotation}deg)`;

        // Stop animation when velocity is very low
        if (Math.abs(rotationVelocity) > 0.01) {
            requestAnimationFrame(animateSpin);
        } else {
            // Snap to the closest section or perform the final action
            snapToSection();
        }
    }

    function snapToSection() {
        // Snap the wheel to the nearest section (adjust the snapping as per your wheel's division)
        const sectionAngle = 360 / 6;  // Assuming 6 sections on the wheel
        const finalRotation = Math.round(currentRotation / sectionAngle) * sectionAngle;
        rotatableImage.style.transition = 'transform 0.5s ease-out';
        rotatableImage.style.transform = `rotate(${finalRotation}deg)`;
        currentRotation = finalRotation;
        setTimeout(() => {
            rotatableImage.style.transition = '';  // Reset the transition after snap
        }, 500);
        
        // Here you can also determine which section is selected and trigger corresponding actions
        let selectedIndex = (Math.round(currentRotation / sectionAngle) + 6) % 6;
        console.log("Selected section:", selectedIndex);
    }

    rotatableImage.addEventListener('mouseleave', function() {
        if (isDragging) {
            isDragging = false;
            rotatableImage.style.cursor = 'grab';
        }
    });
});
