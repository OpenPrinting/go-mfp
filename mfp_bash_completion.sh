#!/bin/bash

# mfp-xxx bash autocompletion script.
#
# put this file in the /etc/bash_completion.d/ or source it from your
# $HOME/.bashrc

__mfp_complete_log()
{
    # Uncomment for logging
    #echo "$@" >> __mfp_complete.log
    return
}

__mfp_complete()
{
	local	cmd args words cword cur cur_raw reply r r2 prefix

	__mfp_complete_log =====

	# Shell passes us command line split into words in the
	# COMP_WORDS array, but its idea about word separator
	# characters doesn't match our goals. In particular, it
	# considers '=' character the word separator, which breaks
	# the long options and considers ':' character the word
	# separator, which breaks the host:port parameters and
	# HTTP URLs.
	#
	# _get_comp_words_by_ref from the bash-completion package re-splits
	# the command line allowing some characters to be excluded
	# from the word separators list with the -n option. Nut now
	# we have another problem: shell automatically removes current
	# word prefix from the COMPREPLY strings, but after re-split the
	# current word might have been extended with the additional
	# prefix which shell doesn't know about and hence cannot
	# remove, so we have to help it.
	_get_comp_words_by_ref -n "=:" words cword cur
	cur_raw="${COMP_WORDS[COMP_CWORD]}"

	if [ "${cur}" != "${cur_raw}" ]; then
		prefix="${cur}"
		prefix="${cur%"${cur_raw}"}"

		if [ "${cur_raw}" == "=" ]; then
			prefix="${prefix}${cur_raw}"
		fi
	fi

	cmd="$1"
	args=("${words[@]:1:$cword}")
	IFS=$'\n' read -r -d '' -a reply < <("$cmd" --bash-completion "${args[@]}" && printf '\0')

	__mfp_complete_log "cur: $cur, cur_raw: $cur_raw, prefix: $prefix"

	# Build COMPREPLY, removing prefix
	COMPREPLY=()
	for r in "${reply[@]}"; do
		r2="${r#"${prefix}"}"
		COMPREPLY+=("${r2}")

		__mfp_complete_log strip: "$r->$r2" "(prefix: ${prefix})"
	done
}

complete -o nospace -F __mfp_complete mfp-cups
complete -o nospace -F __mfp_complete mfp-discover
complete -o nospace -F __mfp_complete mfp-model
complete -o nospace -F __mfp_complete mfp-ppd
complete -o nospace -F __mfp_complete mfp-proxy
complete -o nospace -F __mfp_complete mfp-usb
complete -o nospace -F __mfp_complete mfp-virtual
