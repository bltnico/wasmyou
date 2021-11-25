'use strict';

const root = document.documentElement;
const metaEl = document.querySelector('meta[name="theme-color"]');
const imgInputEl = document.getElementById('img-input');
const buttonEl = document.getElementById('button');
const hoursEl = document.getElementById('hours');
const overlayEl = document.getElementById('overlay');
const [hours, minutes] = hoursEl.getElementsByTagName('span');

imgInputEl.addEventListener('change', () => {
  const file = imgInputEl.files[0];
  if (!file) {
    return;
  }

  // 1MB in Bytes is 1,048,576
  if (file.size > 1048576 * 5) {
    alert('Very large image 😬');
    return;
  }

  root.style.setProperty('--alpha', 0);
  root.style.setProperty('--zoom', 5);

  requestAnimationFrame(async () => {
    const [_, colors] = await Promise.all([
      showPreview(file),
      getCommonColors(file),
    ]);

    setColors(colors);
  });

  root.style.setProperty('--alpha', 1);
  root.style.setProperty('--zoom', 1);
});

function addZero(i) {
  if (i < 10) {
    return '0' + i;
  }

  return i;
}

function showPreview(file) {
  return new Promise((resolve) => {
    const sourceImgEl = document.getElementById('source-img');
    const reader = new FileReader();

    reader.onload = function(event) {
      const now = new Date();
      hours.innerHTML = addZero(now.getHours());
      minutes.innerHTML = addZero(now.getMinutes());
      sourceImgEl.setAttribute('src', event.target.result);
      resolve();
    }

    reader.readAsDataURL(file);
  });
}

function getCommonColors(file) {
  return new Promise((resolve) => {
    const reader = new FileReader();

    reader.onload = async (event) => {
      const bytes = new Uint8Array(event.target.result);
      const colors = await findCommonColors(bytes);
      resolve(colors);
    }

    reader.readAsArrayBuffer(file);
  });
}

function setColors(colors) {
  const [firstColor, secondColor, mainColor] = colors;
  root.style.setProperty('--first-color', firstColor);
  root.style.setProperty('--background-color', firstColor);
  root.style.setProperty('--second-color', secondColor);
  root.style.setProperty('--main-color', mainColor);
  root.style.setProperty('--text-color', mainColor);
  metaEl.setAttribute('content',  firstColor);
  overlayEl.style.backgroundColor = firstColor.replace(/[^,]+(?=\))/, '0.3');
}
