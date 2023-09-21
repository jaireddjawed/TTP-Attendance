const engineeringMajorDropdown = document.querySelector("select[name='engineering-major']");

// display the other major input if the student selects other in the major dropdown options
engineeringMajorDropdown.addEventListener("change", (event) => {
  const engineeringMajor = engineeringMajorDropdown.value;

  if (engineeringMajor === 'other') {
    document.querySelector('.other').style.display = 'block';
  } else {
    document.querySelector('.other').style.display = 'none';
  }
});

document.querySelector('form').addEventListener('submit', (event) => {
  event.preventDefault();

});
