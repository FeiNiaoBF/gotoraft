"use client"

import { ArrowLeft } from "lucide-react"
import Link from "next/link"
import { useLanguage } from "@/contexts/language-context"

export function BackToHome() {
  const { t } = useLanguage()

  return (
    <Link
      href="/"
      className="inline-flex items-center gap-2 px-4 py-2 text-sm text-slate-400 hover:text-white transition-colors"
    >
      <ArrowLeft className="w-4 h-4" />
      {t("backToHome")}
    </Link>
  )
}

