import { Button } from "@/components/animate-ui/components/buttons/button";
import EditableAvatar from "@/components/editable-avatar";
import {
	Card,
	CardContent,
	CardFooter,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import { useForm } from "@tanstack/react-form";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/_settings/account")({
	component: RouteComponent,
	staticData: {
		title: "Account Settings",
	},
});

function RouteComponent() {
	const form = useForm({
		defaultValues: {},
	});

	return (
		<form>
			<Card>
				<CardHeader>
					<CardTitle>
						<EditableAvatar alt="User" fallback={"US"} size={64} />
					</CardTitle>
				</CardHeader>
				<CardContent></CardContent>
				<CardFooter>
					<Button>Save Changes</Button>
				</CardFooter>
			</Card>
		</form>
	);
}
