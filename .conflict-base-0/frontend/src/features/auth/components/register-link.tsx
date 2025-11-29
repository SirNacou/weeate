import { Link } from "@tanstack/react-router";

const RegisterLink = () => {
  return (
    <div className="mt-4 text-center text-sm">
      Don&apos;t have an account?{" "}
      <Link to="/sign-up" className="underline underline-offset-4">
        Sign up
      </Link>
    </div>
  );
};

export default RegisterLink;
