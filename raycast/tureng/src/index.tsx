import { ActionPanel, List, Action, Color, LaunchProps } from "@raycast/api";
import { useState, useEffect } from "react";
import { Client, TranslationResult, TermResult } from "./tureng/tureng";

interface Arguments {
    term: string;
}
export default function Command(props: LaunchProps<{ arguments: Arguments }>) {
    const { term } = props.arguments;
    const [translationResult, setTranslationResult] = useState<TranslationResult | null>();


    useEffect(() => {
        const translate = async () => {
            const client = new Client();
            const result: TranslationResult | null = await client.translate(term);
            setTranslationResult(result);
        }

        translate();
    }, [term]);


    const getTermItems = () => {
        if (!translationResult) {
            return []
        }

        const termResults = translationResult.SearchedTerm == translationResult.PrimeATerm ? translationResult.AResults : translationResult.BResults;
        if (!termResults) {
            return [];
        }


        const getTerm = (termResult: TermResult) => {
            return translationResult.SearchedTerm === translationResult.PrimeATerm ? termResult.TermB : termResult.TermA;
        }
        const getCategory = (termResult: TermResult) => {
            return translationResult.SearchedTerm === translationResult.PrimeATerm ? termResult.CategoryTextB : termResult.CategoryTextA;
        }

        const getTermType = (termResult: TermResult) => {
            return translationResult.SearchedTerm === translationResult.PrimeATerm ? termResult.TermTypeTextB : termResult.TermTypeTextA;
        }


        return termResults.map((termResult: TermResult) => {
            const term = getTerm(termResult);
            const category = getCategory(termResult);
            const termType = getTermType(termResult);

            return {
                id: `${term}-${category}-${termType}`,
                term: term,
                category: category,
                termType: termType,
            }
        });
    }



    return (
        <List>
            {getTermItems().map((item) => {
                return (<List.Item
                    key={item.id}
                    title={item.term}
                    actions={
                        <ActionPanel title="Actions">
                            <Action.CopyToClipboard title="Copy to clipboard" content={item.term} />
                        </ActionPanel>
                    }
                    accessories={[
                        { tag: { value: item.category, color: Color.PrimaryText } },
                        { tag: { value: item.termType, color: Color.SecondaryText } },
                    ]}
                />)
            })}
        </List>
    );
}

