# Graph Report - app-listener  (2026-10-14)

## Corpus Check
- 370 files · ~1,565,458 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4680 nodes · 19485 edges · 177 communities (150 shown, 27 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1644 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `6d69dccc`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- Candidate
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- globreserve_test.go
- filevault.go
- LibraryClosure
- NetworkGuardUseCase
- go_pkg_path_filepath
- install.go
- IntegrationSuite
- IntegrationSuite
- guard_trust.bpf.c
- time.Duration
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- net_tester/main.go
- .FilterExistingWhitelist
- guard.bpf.c
- NetworkMonitor
- IntegrationSuite
- KernelDev
- hashCache
- backups.go
- Serve
- guardSpec
- catalogRefresher
- fileedit_model_test.go
- model
- offline.go
- os.FileMode
- sync.Mutex
- BinaryEntry
- Finder
- main_test.go
- absPath
- Monitor
- runLockdown
- update_test.go
- github.com/charmbracelet/bubbles/textarea.Model
- golang.org/x/sys/unix.Stat_t
- IntegrationSuite
- planUpdaters
- inode_dev
- netModel
- daemonconfig.go
- github.com/charmbracelet/bubbletea.Cmd
- TrustGuard
- testKeyPair
- inodeChain
- GuardInodeKey
- changelog_test.go
- releaseAssets
- fileedit_symlink_test.go
- exe_supersede.h
- GuardEvent
- string
- pinstate.go
- bpfstats/main.go
- runNetworkGuard
- changelog.go
- guarded_ancestor_within_limit
- fcntl
- copyRegularAt
- IntegrationSuite
- time.Time
- btrfsLayoutFrom
- prepareDaemonStart
- github.com/spf13/cobra.Command
- fakeDaemon
- process_vm_readv.c
- downloadAndVerify
- io.Writer
- sanitizeTerminalText
- is_event_type_allowed
- Open
- Ledger
- render.go
- Run
- io.Reader
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- launchKey
- DaemonEvent
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- ldConfParser
- ResolvePinBase
- IntegrationSuite
- Well-known bypass classes
- bootstrapSelfProtectPlaceholders
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- CandidateDir
- Build & sign release assets (reusable workflow)
- FileEvent
- eventUnitTest
- .runEditorHarness
- app-listener Terminal Demo Recording
- IntegrationSuite
- .refreshOnce
- configedit_test.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- parseGuardEvent
- NewDumpHook
- WebSocket /ws client
- stdio
- install.sh
- mc_tag.h
- writeWithin
- make test-integration (rootful Docker suite)
- supersede_probe.c
- trace-app-libs.sh
- IntegrationSuite
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- fakeGuardRepo
- launch_rule_at
- Config
- Guard
- bottomBar
- newTestVetter
- fscrypt.go
- daemon: permanent protection with encryption at rest
- ConfinedEntry
- highlight_test.go
- tempGrantSetup
- Installation
- os.File
- graphify-refresh.sh
- RawTarget
- netGuardModel
- formatEntry
- Vault
- Limitations
- Troubleshooting
- Daemon configuration
- edit-protected
- IntegrationSuite
- Compatibility
- How it works
- buildOneGuard
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- monitor: see what touches a path
- kernel-prof.sh
- vaultFS
- openBunTmp

## God Nodes (most connected - your core abstractions)
1. `Guard` - 114 edges
2. `Config` - 95 edges
3. `GuardInodeKey` - 94 edges
4. `Load()` - 89 edges
5. `IntegrationSuite` - 81 edges
6. `absPath()` - 79 edges
7. `writeConfig()` - 69 edges
8. `IntegrationSuite` - 66 edges
9. `fileEditModel` - 60 edges
10. `engine` - 58 edges

## Surprising Connections (you probably didn't know these)
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `Well-known bypass classes` --references--> `btrfs_search POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `Well-known bypass classes` --references--> `io_uring POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Shared bubbletea TUI layout across monitor and guard modes** — media_demo_monitor_tui, media_demo_guard_whitelist_tui, media_demo_guard_blacklist_tui, media_demo_tui_resource_footer [INFERRED 0.85]
- **Reserved-name trust defenses (trust object #3)** — integrationtests_exploits_readme_glob_plant, integrationtests_exploits_readme_lib_probe, readme_bun_tmpdir_redirect, claude_catalog_refresh [INFERRED 0.85]
- **Release build, sign and verify pipeline** — _github_workflows_release_prerelease_workflow, _github_workflows_promote_release_promote_release_workflow, _github_workflows__build_release_build_release_workflow, _github_workflows__build_release_ed25519_signing, readme_install_sh, readme_update [INFERRED 0.95]
- **Verifier budget measurement and gating** — _github_workflows_verifier_verifier_gate, claude_bpfstats, readme_daemon_check, claude_verifier_complexity_budget, _github_workflows__build_release_ebpf_regeneration [INFERRED 0.95]

## Communities (177 total, 27 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.11
Nodes (27): Execute(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, multicallCacheKey (+19 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (91): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestDiffCatalogFlagRegistered(), TestIsEncryptedRegularFile(), TestFscryptHarness() (+83 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote(), watchPathsInConfig() (+2 more)

### Community 3 - "Candidate"
Cohesion: 0.04
Nodes (88): logRefreshChange(), TestLogRefreshChange(), appendSectionsAndEdit(), collectDiffAdditions(), configPaths(), pathCovered(), TestDiffMergePreservesExistingAndParses(), TestPathCovered() (+80 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (78): TestPatchCatalogSectionGroupedConfig(), TestSetSectionWhitelistPreservesGroupStructure(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines() (+70 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (38): admit_refusal(), check_watched(), code_vouched(), emit_gate(), env_loader(), env_step(), exe_inode_of(), fill_current_image() (+30 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (31): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+23 more)

### Community 7 - "shellrc.go"
Cohesion: 0.09
Nodes (55): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir() (+47 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.07
Nodes (19): addErrKind, guardUnitTest, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), runTool() (+11 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.17
Nodes (35): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+27 more)

### Community 10 - "filevault.go"
Cohesion: 0.11
Nodes (40): harnessSuite, classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD() (+32 more)

### Community 11 - "LibraryClosure"
Cohesion: 0.09
Nodes (24): defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), LdSoPreloadPaths(), LibraryClosure(), mountVouchesOwnership() (+16 more)

### Community 12 - "NetworkGuardUseCase"
Cohesion: 0.13
Nodes (14): NetworkGuardRepository, NewGuardUseCase(), NetworkGuardUseCase, NewNetworkGuardUseCase(), NewNetworkMonitorUseCase(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle() (+6 more)

### Community 13 - "go_pkg_path_filepath"
Cohesion: 0.32
Nodes (3): isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 14 - "install.go"
Cohesion: 0.05
Nodes (76): reloadDaemonForPasswordChange(), applyDiffAdditions(), restoreDaemonAfterDiffAbort(), runDiffCatalog(), applyLiveRefresh(), buildBinaryIfNeeded(), checkRunningBinaryMatchesInstalled(), cleanupBackups() (+68 more)

### Community 15 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (7): eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (50): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+42 more)

### Community 18 - "time.Duration"
Cohesion: 0.06
Nodes (28): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), holdForward() (+20 more)

### Community 19 - "fileEditModel"
Cohesion: 0.11
Nodes (11): clipLabel(), humanSize(), isBinaryContent(), listEntries(), sortNodes(), TestIsBinaryContent(), chownUser, fileEditMode (+3 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.31
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.16
Nodes (45): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+37 more)

### Community 22 - "Resource"
Cohesion: 0.08
Nodes (12): ownerPaths(), Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots(), TestPartitionEncryptionRootsSplitsDirsAndFiles() (+4 more)

### Community 23 - "net_tester/main.go"
Cohesion: 0.06
Nodes (43): ParseEventsFlag(), Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode() (+35 more)

### Community 24 - ".FilterExistingWhitelist"
Cohesion: 0.11
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.11
Nodes (52): exe_is_superseded(), exe_refused(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade() (+44 more)

### Community 26 - "NetworkMonitor"
Cohesion: 0.08
Nodes (15): NetEvent, NetEventType, NetEventTypes(), ParseNetEventType(), eventsetSummary(), modeLabel(), NewNetGuard(), NewNetworkMonitor() (+7 more)

### Community 28 - "KernelDev"
Cohesion: 0.08
Nodes (35): fdStat, BinaryStat, openRoot(), physicalParent(), statFD(), checkBtrfsKeys(), isFuseType(), parseMajorMinor() (+27 more)

### Community 29 - "hashCache"
Cohesion: 0.18
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 30 - "backups.go"
Cohesion: 0.07
Nodes (40): deletePostBackups(), restoreBackups(), pickUsers(), offerBackupCleanup(), runUninstall(), detectSSHAgentUnits(), removeKeyAndEmptyDir(), removeMasterKey() (+32 more)

### Community 31 - "Serve"
Cohesion: 0.08
Nodes (29): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+21 more)

### Community 32 - "guardSpec"
Cohesion: 0.20
Nodes (12): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+4 more)

### Community 33 - "catalogRefresher"
Cohesion: 0.10
Nodes (9): drain(), newCatalogRefresher(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), catalogRefresher, debouncer (+1 more)

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.15
Nodes (28): OpenRegularNoFollow(), fileModeToUnixMode(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd() (+20 more)

### Community 35 - "model"
Cohesion: 0.22
Nodes (7): formatType(), listenForEvents(), syncViewport(), tickStats(), eventLine, eventMsg, model

### Community 36 - "offline.go"
Cohesion: 0.08
Nodes (40): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), collectFilesystemPrereqs(), confirmRunPrereq(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles(), cleanOrphanedMetadata() (+32 more)

### Community 37 - "os.FileMode"
Cohesion: 0.08
Nodes (35): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig(), runEditConfig(), updateCatalogConfig() (+27 more)

### Community 38 - "sync.Mutex"
Cohesion: 0.08
Nodes (23): readConfigPut(), TestReadGrantRequestConfig(), newControlManager(), readAuthRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines() (+15 more)

### Community 39 - "BinaryEntry"
Cohesion: 0.06
Nodes (28): binaryVerifyState, deferredBinary, HashStore, sharedHashEntry, BinaryEntry, TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), BinariesSummary() (+20 more)

### Community 40 - "Finder"
Cohesion: 0.16
Nodes (19): findButton, span, findAll(), Finder, lowerRunes(), NewFinder(), newFindInput(), overlay() (+11 more)

### Community 41 - "main_test.go"
Cohesion: 0.05
Nodes (20): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+12 more)

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (4): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 43 - "Monitor"
Cohesion: 0.06
Nodes (11): evalSymlinksOrEmpty(), NewMonitor(), newPathCache(), dirTarget(), TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace(), Monitor (+3 more)

### Community 44 - "runLockdown"
Cohesion: 0.20
Nodes (14): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), runLockdown(), recoverPinState(), readTempJournal(), stripTempJournal(), TestTempJournalRoundTripInPlace() (+6 more)

### Community 45 - "update_test.go"
Cohesion: 0.10
Nodes (10): fetchReleases(), readInstalledVersion(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel(), TestNewerThanInstalled(), TestNewerThanStable(), TestParseStableVersion() (+2 more)

### Community 46 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.11
Nodes (33): class, viewportState, classOf(), column(), dedent(), Edit(), leadingIndent(), lines() (+25 more)

### Community 48 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 50 - "planUpdaters"
Cohesion: 0.15
Nodes (17): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+9 more)

### Community 51 - "inode_dev"
Cohesion: 0.24
Nodes (8): ctx_ptr(), inode_dev(), sb_dev(), guard_exec_applet(), guard_inode_free(), mc_basename(), guarded_root_kind(), trust_inode_free()

### Community 52 - "netModel"
Cohesion: 0.13
Nodes (16): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+8 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.11
Nodes (44): systemPatterns(), systemRule(), grantNewBinaries(), watchGroup, Line(), Lines(), TestLinesKeyOnLinkAndIncludePending(), addResource() (+36 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.08
Nodes (11): diffModel, editorModel, clamp(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), noticeModel, servedModel (+3 more)

### Community 55 - "TrustGuard"
Cohesion: 0.05
Nodes (24): tempGrant, trustHook, trustRows, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), dropSupersededLocked(), PruneSuperseded(), setSupersedeTrust() (+16 more)

### Community 56 - "testKeyPair"
Cohesion: 0.23
Nodes (10): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), testKeyPair(), TestParseChecksum(), TestParsePublicKey(), TestVerifyRelease() (+2 more)

### Community 57 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.04
Nodes (42): engine, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates() (+34 more)

### Community 59 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 60 - "releaseAssets"
Cohesion: 0.22
Nodes (10): assetsFor(), checkUpdatePreconditions(), resolveAssets(), TestAssetsFor(), TestResolveAssetsPicksArch(), TestValidateChannel(), validateChannel(), runsNatively() (+2 more)

### Community 61 - "fileedit_symlink_test.go"
Cohesion: 0.16
Nodes (16): AtomicWriteAt(), openRootT(), TestWriteFileKeepMeta(), TestWriteFileKeepMetaRefusesSymlink(), mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent() (+8 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.19
Nodes (17): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+9 more)

### Community 63 - "GuardEvent"
Cohesion: 0.07
Nodes (22): displayEntries(), holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless(), runGuardTUI(), TestDispatchFullChannelNeverDropsDenials() (+14 more)

### Community 64 - "string"
Cohesion: 0.14
Nodes (5): dump_via_debugfs(), main(), probe_open(), denied(), main()

### Community 65 - "pinstate.go"
Cohesion: 0.26
Nodes (14): pinOwnerLikelyAlive(), readPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID(), TestPinStateRoundTrip(), TestReadPinStateMalformed(), TestRecoverPinStateNoRecord() (+6 more)

### Community 66 - "bpfstats/main.go"
Cohesion: 0.38
Nodes (9): chownToSudoUser(), compare(), dumpFullLog(), dumpTrace(), loadOne(), main(), measure(), resizeMaps() (+1 more)

### Community 67 - "runNetworkGuard"
Cohesion: 0.11
Nodes (26): CheckEBPF(), ParseNetEventsFlag(), runMonitor(), runTUI(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode() (+18 more)

### Community 68 - "changelog.go"
Cohesion: 0.22
Nodes (19): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), showChangelog(), compareStableVersions(), confirmUpdate(), filterChannel() (+11 more)

### Community 69 - "guarded_ancestor_within_limit"
Cohesion: 0.27
Nodes (16): evict_inode_from_guard(), get_dentry_from_path(), get_inode_from_path(), guard_path_link(), guard_path_mkdir(), guard_path_mknod(), guard_path_rename(), guard_path_rmdir() (+8 more)

### Community 70 - "fcntl"
Cohesion: 0.16
Nodes (4): main(), report(), denied(), main()

### Community 71 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 73 - "time.Time"
Cohesion: 0.08
Nodes (25): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), markReady(), notifySystemdReady() (+17 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "prepareDaemonStart"
Cohesion: 0.20
Nodes (10): CheckBPFLSM(), configureDaemonLogging(), loadDaemonConfig(), newPinGeneration(), prepareDaemonStart(), resolveConfigPath(), runBPFCheck(), runOneShotMode() (+2 more)

### Community 76 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (19): AddServeFlags(), credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+11 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "downloadAndVerify"
Cohesion: 0.27
Nodes (8): downloadAndVerify(), downloadFile(), downloadReleaseFiles(), sanityCheckBinary(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestSanityCheckBinaryRefusesOtherMachine(), updateHTTPClient()

### Community 81 - "io.Writer"
Cohesion: 0.43
Nodes (7): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), candidate

### Community 82 - "sanitizeTerminalText"
Cohesion: 0.17
Nodes (8): formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision(), formatGuardType(), clipToWidth(), formatResourceBar(), TestSanitizeTerminalText()

### Community 83 - "is_event_type_allowed"
Cohesion: 0.44
Nodes (10): emit_event(), guard_net_socket_bind(), guard_net_socket_connect(), guard_net_socket_listen(), guard_net_socket_recvmsg(), guard_net_socket_sendmsg(), is_event_type_allowed(), is_socket_guarded() (+2 more)

### Community 84 - "Open"
Cohesion: 0.29
Nodes (8): EnsurePlaceholders(), JournalPath(), Open(), TestEnsurePlaceholdersNoInstallDir(), TestLedgerKeepsJournalFile(), TestLedgerRecordLookupPending(), TestHashCacheRoundTripAndPrune(), TestHashCacheSkipsMalformedRow()

### Community 85 - "Ledger"
Cohesion: 0.15
Nodes (3): Ledger, Pending, Source

### Community 86 - "render.go"
Cohesion: 0.16
Nodes (6): lineStyles, lineLen(), markStyle(), repeatSpaces(), stylesOf(), wrap()

### Community 87 - "Run"
Cohesion: 0.07
Nodes (14): columnState, dataRow, dataRowRenderer, guiModel, headerRenderer, headerWidget, findTopLevelWindow(), floatWindow() (+6 more)

### Community 88 - "io.Reader"
Cohesion: 0.33
Nodes (4): BtrfsMounted(), mountinfoHasFstype(), TestMountinfoHasFstype(), progressReader

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 92 - "DaemonEvent"
Cohesion: 0.07
Nodes (45): UIDResolver, NewUIDResolver(), awaitStartupOrSignal(), catchLifecycleSignals(), newDaemonModel(), runDaemon(), runDaemonUI(), runHeadless() (+37 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.07
Nodes (20): reasonText(), eventTypeKey(), mcNameKey(), statInodeKey(), codeVouched(), isTraced(), judgePasses(), loaderEnv() (+12 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "ldConfParser"
Cohesion: 0.50
Nodes (3): ldConfParser, ldConfInclude(), ldSoConfDirs()

### Community 100 - "ResolvePinBase"
Cohesion: 0.50
Nodes (5): dirIsRootOwnedSafe(), isBpffs(), mountBpffs(), ResolvePinBase(), TestDirIsRootOwnedSafe()

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "bootstrapSelfProtectPlaceholders"
Cohesion: 0.67
Nodes (4): ensurePinStateFilePlaceholder(), ensurePlaceholder(), bootstrapSelfProtectPlaceholders(), ensureHashFilePlaceholder()

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "CandidateDir"
Cohesion: 0.05
Nodes (33): entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, GlobChild, statFunc, compileGlobNames() (+25 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "FileEvent"
Cohesion: 0.12
Nodes (9): runHeadless(), FileEvent, Run(), MonitorRepository, MonitorUseCase, NewMonitorUseCase(), newFakeMonitorRepo(), TestMonitorUseCaseStartError() (+1 more)

### Community 110 - "eventUnitTest"
Cohesion: 0.18
Nodes (7): BpfEvent, eventUnitTest, NetBpfEvent, Cstr(), DecodeBpfEvent(), FormatAddr(), Ntohs()

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 115 - ".refreshOnce"
Cohesion: 0.21
Nodes (11): sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), RefreshCatalog(), TestRefreshCatalogAdmitVetsOnlyNewLines() (+3 more)

### Community 116 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "parseGuardEvent"
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 121 - "NewDumpHook"
Cohesion: 0.20
Nodes (12): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+4 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "stdio"
Cohesion: 0.12
Nodes (6): copy(), is_verb(), main(), usage(), put(), wait_then_read()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "fakeGuardRepo"
Cohesion: 0.04
Nodes (27): EventType, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), SupersededBinary, binaryStatOf(), TempBinary, TemporaryGrant, TempRule (+19 more)

### Community 143 - "launch_rule_at"
Cohesion: 0.50
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 144 - "Config"
Cohesion: 0.08
Nodes (41): makeReloadHandler(), noteReloadedSettings(), startCatalogRefresh(), inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet() (+33 more)

### Community 146 - "Guard"
Cohesion: 0.07
Nodes (11): backingDeviceUnion(), exeRow, rootHandle, tempMaskOp, tempRow, SharedPinDegraded(), BackingDevices(), deleteResKeys() (+3 more)

### Community 148 - "bottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 149 - "newTestVetter"
Cohesion: 0.44
Nodes (10): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+2 more)

### Community 150 - "fscrypt.go"
Cohesion: 0.05
Nodes (49): confirmMasterKeyOverwrite(), runGenKey(), askDecrypt(), decryptDirectories(), decryptStep(), pickDirsToDecrypt(), TestDecryptDirectoriesEmpty(), deprovisionKind (+41 more)

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "ConfinedEntry"
Cohesion: 0.29
Nodes (19): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+11 more)

### Community 156 - "highlight_test.go"
Cohesion: 0.13
Nodes (20): nextConfig(), fg(), Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo() (+12 more)

### Community 157 - "tempGrantSetup"
Cohesion: 0.57
Nodes (7): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure()

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "os.File"
Cohesion: 0.05
Nodes (54): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), setConfinedHomes(), binaryVetter, vetDecision, VettedInode, Multicall (+46 more)

### Community 161 - "RawTarget"
Cohesion: 0.24
Nodes (11): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+3 more)

### Community 164 - "netGuardModel"
Cohesion: 0.24
Nodes (5): NetGuardEvent, protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 168 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 185 - "buildOneGuard"
Cohesion: 0.08
Nodes (40): buildConcurrency(), buildGuards(), buildOneGuard(), drainEphemeral(), eventFilterOptions(), reloadOnce(), reloadSlotsNeeded(), startGuardedDaemon() (+32 more)

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.40
Nodes (5): app-listener documentation, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network-monitor: see what a program talks to

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "vaultFS"
Cohesion: 0.18
Nodes (9): checkFreshDir(), openDirNoFollow(), renewDir(), TestRuntimeClassMarkers(), applyNewFileMeta(), defaultNewFileMode(), TestDefaultNewFileMode(), writeFileInPlace() (+1 more)

### Community 201 - "openBunTmp"
Cohesion: 0.39
Nodes (8): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp, inodeID

## Knowledge Gaps
- **32 isolated node(s):** `IntegrationSuite`, `github.com/Virgula0/app-listener`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 449 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **27 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Process inspectors` connect `time.Time` to `Daemon configuration`?**
  _High betweenness centrality (0.049) - this node is a cross-community bridge._
- **What connects `IntegrationSuite`, `github.com/Virgula0/app-listener`, `open_by_handle_at POC` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10934014986341879 - nodes in this community are weakly interconnected._
- **Why does `Daemon configuration` connect `Daemon configuration` to `time.Time`, `README.md`?**
  _High betweenness centrality (0.049) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.028368794326241134 - nodes in this community are weakly interconnected._
- **Why does `watchPattern` connect `time.Time` to `globreserve_test.go`, `go_pkg_path_filepath`, `catalogRefresher`?**
  _High betweenness centrality (0.048) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.0898838004101162 - nodes in this community are weakly interconnected._