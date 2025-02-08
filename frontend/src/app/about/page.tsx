"use client"

import { PageLayout } from "@/components/page-layout"
import { BackToHome } from "@/components/back-to-home"
import { useLanguage } from "@/contexts/language-context"

export default function AboutPage() {
  const { t } = useLanguage()

  return (
    <PageLayout
      title={t("aboutTitle")}
      description={t("Learn about the team behind this Raft algorithm visualization tool.")}
    >
      <div className="mb-6">
        <BackToHome />
      </div>
      <div className="bg-slate-800/50 p-6 rounded-lg">
        <p className="text-slate-300">{t("Project and team information will be implemented here.")}</p>
      </div>
    </PageLayout>
  )
}

