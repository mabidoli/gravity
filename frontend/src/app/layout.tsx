import type { Metadata } from "next";
import { ClerkProvider } from "@clerk/nextjs";
import { dark } from "@clerk/themes";
import { Providers } from "./providers";
import "./globals.css";

export const metadata: Metadata = {
  title: "Gravity V2 - Personal Infrastructure",
  description: "Unified communication, calendar, and social streams in one workspace",
  keywords: ["productivity", "unified inbox", "communication hub"],
};

// Check if Clerk is configured
const isClerkConfigured = !!process.env.NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY;

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  // When Clerk is not configured, render without ClerkProvider
  if (!isClerkConfigured) {
    return (
      <html lang="en" className="dark">
        <body className="font-sans antialiased">
          <Providers>{children}</Providers>
        </body>
      </html>
    );
  }

  // When Clerk is configured, wrap Providers inside ClerkProvider
  return (
    <ClerkProvider
      appearance={{
        baseTheme: dark,
        variables: {
          colorPrimary: "#14b8a6",
          colorBackground: "#0f172a",
          colorInputBackground: "#1e293b",
          colorInputText: "#f8fafc",
        },
      }}
    >
      <html lang="en" className="dark">
        <body className="font-sans antialiased">
          <Providers>{children}</Providers>
        </body>
      </html>
    </ClerkProvider>
  );
}
