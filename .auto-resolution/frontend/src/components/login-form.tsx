import { cn } from "@/lib/utils";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useNavigate } from "@tanstack/react-router";
import LoginWithPasswordCard from "@/features/auth/components/login-with-password-card";
import { SeparatorWithText } from "./separator-with-text";
import { LoginWithGoogleButton } from "@/features/auth/components/login-with-google-button";
import RegisterLink from "@/features/auth/components/register-link";

export function LoginForm({
  className,
  ...props
}: React.ComponentPropsWithoutRef<"div">) {
  const navigate = useNavigate();

  async function onLoginSuccess(): Promise<void> {
    await navigate({ to: "/protected" });
  }

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader>
          <CardTitle className="text-2xl">Login</CardTitle>
          <CardDescription>
            Enter your email below to login to your account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <LoginWithPasswordCard
            onSuccess={onLoginSuccess}
            onError={function (error: string): void {
              throw new Error("Function not implemented.");
            }}
          />
          <SeparatorWithText text="Or Login With" />
          <LoginWithGoogleButton />
          <RegisterLink />
        </CardContent>
      </Card>
    </div>
  );
}
