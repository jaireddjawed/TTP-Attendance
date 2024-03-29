// todo: request permission to use full screen mode

/*
  * This function submits the sign in information to the server
  * and redirects the user to the success page if the sign in was successful.
*/

function submitSignInInformation(firstName, lastName, studentId) {
  fetch('/submit-signin', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      'first-name': firstName,
      'last-name': lastName,
      'student-id': studentId,
    })
  })
  .then(response => {
    if (!response.ok) {
      alert('There was an error signing you in. Please try again.');
      return;
    }

    return response.text();
  })
  .then(message => {
    // if statement just for Marisol 😂
    if (firstName === 'MARISO') {
      firstName = 'MARISOL';
    }

    // redirect to success page or major page depending on whether the student has already signed in before
    if (message === 'success') {
      window.location.href = '/success?firstName=' + firstName + '&lastName=' + lastName;
    }
    else if (message === 'student-not-exist') {
      window.location.href = '/major?firstName=' + firstName + '&lastName=' + lastName + '&studentId=' + studentId;
    }
  });
}

/*
 * form validation to ensure that everyone
 * fills out every field before submitting
 * the sign in form
 */

$(document).ready(function() {
  $('.form').validate({
    rules: {
      'first-name': {
        required: true,
      },
      'last-name': {
        required: true,
      },
      'student-id': {
        required: true,
        minlength: 9,
        maxlength: 9,
        digits: true,
      },
    },
    messages: {
      'first-name': {
        required: 'Please enter your first name.',
      },
      'last-name': {
        required: 'Please enter your last name.',
      },
      'student-id': {
        required: 'Please enter your student id.',
        minlength: 'Your student id number must be 9 digits long.',
        maxlength: 'Your student id number must be 9 digits long.',
        digits: 'Your student id number can only contain digits.',
      },
    },
    submitHandler: function(form) {
      submitSignInInformation(
        document.querySelector('input[name="first-name"]').value,
        document.querySelector('input[name="last-name"]').value,
        document.querySelector('input[name="student-id"]').value,
      );
    },
  });
});

/*
 * This event listener is for the id card reader. Once a student swipes their id card, it parses their
 * first and last name from the id card string and submits the sign in form.
 */
document.querySelector('input[name="id-card-reader"]').addEventListener('change', (event) => {
  let fullName, studentId;

  try {
    const idString = document.querySelector('input[name="id-card-reader"]').value;
    fullName = idString.split('^')[1].split('/');

    studentId = idString.split('^')[2]
    studentId = studentId.substring(12, 21)

  } catch (error) {
    return;
  }
  const [lastName, firstName] = fullName;
  submitSignInInformation(firstName, lastName, studentId);
});

/*
 * Inactivity needs to be measured in case someone clicks on one of the inputs and
 * doesn't sign in so that the browser can refocus on the id-card input in case
 * the next student signs in with their id card.
 * 
 * If it does not refocus on the id card input, no one will be able to sign in with their student id.
 */

let inactivity = 0;

addEventListener('keypress', () => {
  inactivity = 0;
});

addEventListener('mousemove', () => {
  inactivity = 0;
});

addEventListener('click', () => {
  inactivity = 0;
});

// measure inactivity every second
setInterval(() => {
  inactivity++;

  // refocus after 15 seconds of inactivity`
  if (inactivity >= 15) {
    const inputs = document.querySelectorAll('input');
    inputs.forEach(input => input.value = '');
    document.querySelector('input[name="id-card-reader"]').focus();
  }
}, 1000);
