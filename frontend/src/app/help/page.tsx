"use client"

import { PageLayout } from "@/components/page-layout"
import { BackToHome } from "@/components/back-to-home"
import { useLanguage } from "@/contexts/language-context"

export default function HelpPage() {
  const { t } = useLanguage()

  return (
    <PageLayout
      title={t("helpTitle")}
      description={t("Learn more about the Raft algorithm and how to use this visualization tool.")}
    >
      <div className="mb-6">
        <BackToHome />
      </div>
      <div className="bg-slate-800/50 p-6 rounded-lg">
        <p className="text-slate-300">{t("Help content and documentation will be implemented here.")}</p>
      </div>
    </PageLayout>
  )
}
