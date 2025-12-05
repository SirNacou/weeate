import { useNavigate } from "@tanstack/react-router";
import { toast } from "sonner";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import { LoginWithGoogleButton } from "@/features/auth/components/login-with-google-button";
import LoginWithPasswordCard from "@/features/auth/components/login-with-password-card";
import RegisterLink from "@/features/auth/components/register-link";
import { cn } from "@/lib/utils";
import { SeparatorWithText } from "./separator-with-text";

export function LoginForm({
	className,
	...props
}: React.ComponentPropsWithoutRef<"div">) {
	const navigate = useNavigate();

	async function onLoginSuccess(): Promise<void> {
		await navigate({ to: "/" });
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
						onError={(error: string): void => {
							toast.error(error);
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
