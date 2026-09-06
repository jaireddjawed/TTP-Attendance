/**
 * Parsing for the magnetic stripe on a UCR student id card.
 *
 * The reader behaves as a keyboard: it types the whole track into the focused
 * input and presses Enter. The track is caret-delimited, with the cardholder
 * name in the second field as "LAST/FIRST" and the student id embedded in the
 * third field.
 */

export interface SwipedCard {
  firstName: string;
  lastName: string;
  studentId: string;
}

const ID_OFFSET = 12;
const ID_LENGTH = 9;

/**
 * parseIdCard pulls the name and student id out of a swiped track, returning
 * null for anything that is not a readable card — a mis-swipe, or a person
 * typing into the hidden field by accident.
 */
export function parseIdCard(track: string): SwipedCard | null {
  const fields = track.split("^");
  if (fields.length < 3) {
    return null;
  }

  const [lastName, firstName] = fields[1].split("/");
  if (!firstName || !lastName) {
    return null;
  }

  const studentId = fields[2].substring(ID_OFFSET, ID_OFFSET + ID_LENGTH);
  if (!/^\d{9}$/.test(studentId)) {
    return null;
  }

  return {
    firstName: firstName.trim(),
    lastName: lastName.trim(),
    studentId,
  };
}
