"use client";

import { useEffect } from "react";
import { setTokenGetter } from "@/lib/api";

// Check if Clerk is configured
const isClerkConfigured = !!process.env.NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY;

/**
 * Hook to connect Clerk authentication with the API client.
 * This should be called once at the app level to set up token injection.
 * When Clerk is not configured, it sets up a null token getter.
 */
export function useAuthenticatedApi() {
  useEffect(() => {
    if (!isClerkConfigured) {
      // When Clerk is not configured, set a no-op token getter
      setTokenGetter(async () => null);
    }
  }, []);

  // If Clerk is not configured, return early with default values
  if (!isClerkConfigured) {
    return { isLoaded: true, isSignedIn: false };
  }

  // When Clerk is configured, use the Clerk-specific implementation
  return useAuthenticatedApiWithClerk();
}

/**
 * Internal hook that uses Clerk - only called when Clerk is configured.
 */
function useAuthenticatedApiWithClerk() {
  // eslint-disable-next-line @typescript-eslint/no-require-imports
  const { useAuth } = require("@clerk/nextjs");
  const { getToken, isLoaded, isSignedIn } = useAuth();

  useEffect(() => {
    if (isLoaded) {
      // Set up the token getter for the API client
      setTokenGetter(async () => {
        if (!isSignedIn) {
          return null;
        }
        return getToken();
      });
    }
  }, [getToken, isLoaded, isSignedIn]);

  return { isLoaded, isSignedIn };
}
