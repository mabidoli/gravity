import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Gravity V2 - Personal Infrastructure",
  description: "Unified communication, calendar, and social streams in one workspace",
  keywords: ["productivity", "unified inbox", "communication hub"],
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <body className="font-sans antialiased">
        {children}
      </body>
    </html>
  );
}
