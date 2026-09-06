/** A BCOE major a student can pick the first time they sign in. */
export interface Major {
  value: string;
  label: string;
}

/** OTHER_MAJOR is the option that reveals the free-text major field. */
export const OTHER_MAJOR = "other";

export const MAJORS: Major[] = [
  { value: "bioengineering", label: "Bioengineering" },
  { value: "chemical-engineering", label: "Chemical Engineering" },
  { value: "computer-engineering", label: "Computer Engineering" },
  { value: "computer-science", label: "Computer Science" },
  { value: "computer-science-business-applications", label: "Computer Science w/ Business Applications" },
  { value: "data-science", label: "Data Science" },
  { value: "electrical-engineering", label: "Electrical Engineering" },
  { value: "environmental-engineering", label: "Environmental Engineering" },
  { value: "material-science-engineering", label: "Material Science & Engineering" },
  { value: "mechanical-engineering", label: "Mechanical Engineering" },
  { value: "robotics", label: "Robotics" },
  { value: OTHER_MAJOR, label: "Other" },
];
