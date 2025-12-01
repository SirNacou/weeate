import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { createFileRoute, Link } from "@tanstack/react-router";

export const Route = createFileRoute("/auth/auth-code-error")({
  validateSearch: (search: Record<string, unknown>) => {
    return {
      error: (search.error as string) || "Unknown error occurred",
    };
  },
  component: AuthCodeError,
});

function AuthCodeError() {
  const { error } = Route.useSearch();
  return (
    <div className="flex min-h-screen items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>Authentication Error</CardTitle>
          <CardDescription>There was a problem signing you in</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="rounded-lg bg-destructive/10 p-3 border border-destructive/20">
            <p className="text-sm text-destructive font-medium">{error}</p>
          </div>
          <p className="text-sm text-muted-foreground">
            The authentication code was invalid or has expired. Please try
            signing in again.
          </p>
          <Button asChild className="w-full">
            <Link to="/login">Back to Login</Link>
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
