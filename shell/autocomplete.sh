_gdpx() {
    COMPREPLY=()
    local current="${COMP_WORDS[COMP_CWORD]}"

    local commands=(
        list
        start
        stop
        recreate
        install-autocomplete
    )
    if [[ ${COMP_CWORD} == 1 ]] ; then
        COMPREPLY=( $(compgen -W "${commands[*]}" -- ${current}) )
        return 0
    fi

    local command="${COMP_WORDS[1]}"

    case "${command}" in
        list)
            local subcommands=(
                --started
                --stopped
            )
            COMPREPLY=( $(compgen -W "${subcommands[*]}" -- ${current}) )
            return 0
            ;;
        help)
            COMPREPLY=()
            return 0
            ;;
        recreate)
            COMPREPLY=()
            return 0
            ;;
        install-autocomplete)
            COMPREPLY=()
            return 0
            ;;
        start)
            if [[ ${COMP_CWORD} == 2 ]] ; then
                local stopped_list=$(gdpx list --stopped | awk '{ sub(/^[ \t]+/, ""); print $1}' | sed 's/\x1B\[[0-9;]\{1,\}[A-Za-z]//g' | xargs)
                COMPREPLY=( $(compgen -W "${stopped_list}" -- ${current}) )
                return 0
            fi
            ;;
        stop)
            if [[ ${COMP_CWORD} == 2 ]] ; then
                local started_list=$(gdpx list --started | awk '{ sub(/^[ \t]+/, ""); print $1}' | sed 's/\x1B\[[0-9;]\{1,\}[A-Za-z]//g' | xargs)
                COMPREPLY=( $(compgen -W "${started_list}" -- ${current}) )
                return 0
            fi
            ;;
        esac

}

complete -F _gdpx gdpx
# 