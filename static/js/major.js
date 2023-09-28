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

/*
 * form validation to ensure that everyone
 * fills out every field before submitting
 * the sign in form
 */

$(document).ready(function() {
  $('.form').validate({
    rules: {
      'engineering-major': {
        required: true,
      },
      'other-major': {
        required: true,
      },
    },
    messages: {
      'engineering-major': {
        required: 'Please select a major.',
      },
      'other-major': {
        required: 'Please enter your major.',
      },
    },
    submitHandler: function(form) {

    },
  });
});
