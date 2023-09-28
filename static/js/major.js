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
    submitHandler: function() {
      const params = new URLSearchParams(window.location.search);

      const firstName = params.get('firstName'),
            lastName = params.get('lastName'),
            studentId = params.get('studentId');

      const engineeringMajor = document.querySelector('select[name="engineering-major"]').value,
            otherMajor = document.querySelector('input[name="other-major"]').value;

      fetch('/submit-major', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          'first-name': firstName,
          'last-name': lastName,
          'student-id': studentId,
          'major': engineeringMajor === 'Other' ? otherMajor : engineeringMajor,
        }),
      })
      .then(response => {
        if (!response.ok) {
          alert('There was an error signing you in. Please try again.');
          return;
        }

        window.location.href = `/success?firstName=${firstName}&lastName=${lastName}`;
      });
    },
  });
});
