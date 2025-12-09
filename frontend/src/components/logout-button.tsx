import { useNavigate } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import { createClient } from "@/lib/supabase/client";

export function LogoutButton() {
	const navigate = useNavigate();

	const logout = async () => {
		const supabase = createClient();
		await supabase.auth.signOut();
		navigate({ to: "/login" });
	};

	return <Button onClick={logout}>Logout</Button>;
}
