import { AnyFormState } from "@tanstack/react-form"

export const getFormSubmissionStatus = (state: AnyFormState) => ({
  canSubmit: state.isValid && !state.isPristine && state.canSubmit,
  isSubmitting: state.isSubmitting,
});
