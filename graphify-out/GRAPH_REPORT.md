# Graph Report - app-listener  (2026-10-18)

## Corpus Check
- 471 files · ~1,693,318 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6161 nodes · 26188 edges · 197 communities (165 shown, 32 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2460 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `fe22b15c`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- desktop.go
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- mkdirs
- github.com/charmbracelet/bubbles/textarea.Model
- policyEditor
- icon
- github.com/charmbracelet/bubbletea.Cmd
- .launch
- User
- IntegrationSuite
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- setpassword_authz_test.go
- MountVouch
- guard.bpf.c
- expandPlaceholders
- time.Time
- NewDumpHook
- hashFileShared
- gui/install.go
- serve.go
- IntegrationSuite
- netModel
- picker
- github.com/spf13/cobra.Command
- IntegrationSuite
- backend
- controlServer
- fileedit_model_test.go
- os.File
- IntegrationSuite
- absPath
- monitorUnitTest
- NetGuardEvent
- GuardEvent
- Open
- stdio
- editView
- fadeScroll
- LibraryClosure
- GenerateConf
- proto.go
- daemonconfig.go
- tabs.go
- TrustGuard
- wizard.go
- Network
- GuardInodeKey
- filevault.go
- update/update.go
- newTestTarget
- exe_supersede.h
- NetEventType
- changelog.go
- bpfstats/main.go
- gui/session.go
- runMonitor
- changelog_test.go
- codeedit.go
- fcntl
- runOfflineEdit
- Guard
- sync.Mutex
- btrfsLayoutFrom
- hashCache
- tempGrantSetup
- RefreshCatalog
- prepareDaemonStart
- process_vm_readv.c
- Monitor
- updatePage
- planUpdaters
- dialog_test.go
- .posAt
- find_test.go
- configedit_test.go
- liveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- DaemonEvent
- .FilterExistingWhitelist
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- debouncer
- github.com/stretchr/testify/suite.Suite
- writeEvent
- Well-known bypass classes
- SuperblockDev
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- BinaryStat
- Build & sign release assets (reusable workflow)
- installPage
- mountTable
- TestMonitorUseCaseLifecycle
- IntegrationSuite
- IntegrationSuite
- mountsUnitTest
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- CandidateDir
- IntegrationSuite
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- make test-integration (rootful Docker suite)
- supersede_probe.c
- go_pkg_fyne_io_fyne_v2
- trace-app-libs.sh
- IntegrationSuite
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- fakeGuardRepo
- time.Duration
- IntegrationSuite
- ParseEventsFlag
- Config
- net_tester/main.go
- Execute
- fscrypt.go
- ConfinedEntry
- EventType
- layout.go
- copyRegularAt
- journal
- Installation
- backend_test.go
- graphify-refresh.sh
- vaultfs.go
- inodeChain
- pinstate.go
- feed.go
- shell
- runView
- writeWithin
- fileedit_symlink_test.go
- Limitations
- Troubleshooting
- IntegrationSuite
- Daemon configuration
- KernelDev
- install/install.go
- Compatibility
- How it works
- fakeDaemon
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- NetworkMonitor
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- edit-protected
- monitor: see what touches a path
- kernel-prof.sh
- pathCache
- launchKey
- runNetworkGuard
- Desktop GUI
- IntegrationSuite
- Ledger
- inode_dev.h
- .eventTypeAllowed
- binaryVetter
- bleedLayout
- TestMain

## God Nodes (most connected - your core abstractions)
1. `Guard` - 117 edges
2. `GuardInodeKey` - 96 edges
3. `Config` - 95 edges
4. `Load()` - 89 edges
5. `shell` - 87 edges
6. `IntegrationSuite` - 82 edges
7. `absPath()` - 81 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 69 edges
10. `IntegrationSuite` - 68 edges

## Surprising Connections (you probably didn't know these)
- `Editor keys` --references--> `Delete()`  [INFERRED]
  docs/edit-protected.md → internal/backups/backups.go
- `Demo GIFs` --references--> `TestDemoReel()`  [INFERRED]
  docs/development.md → internal/gui/demo_test.go
- `Demo GIFs` --references--> `TestSnapshots()`  [INFERRED]
  docs/development.md → internal/gui/snapshot_test.go
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md

## Import Cycles
- None detected.

## Communities (197 total, 32 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.12
Nodes (25): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, multicallCacheKey, attrOpCase (+17 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (112): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), awaitStartupOrSignal(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestAwaitStartupOrSignalAbortsAndLocksDown(), TestAwaitStartupOrSignalAbortWithFailedStartup() (+104 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "desktop.go"
Cohesion: 0.32
Nodes (3): isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (79): TestPatchCatalogSectionGroupedConfig(), TestSetSectionWhitelistPreservesGroupStructure(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines() (+71 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (32): emit_event(), emit_event_kern(), fill_identity(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice() (+24 more)

### Community 7 - "shellrc.go"
Cohesion: 0.12
Nodes (41): rcTarget, invokingIDs(), writeUserFile(), BunTmpDir(), bunTargets(), createOwned(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv() (+33 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (22): addErrKind, ReplacementCheck, BinariesSummary(), classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest (+14 more)

### Community 9 - "mkdirs"
Cohesion: 0.13
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 10 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.10
Nodes (27): findButton, pos, span, column(), dedent(), lines(), moveTo(), Navigate() (+19 more)

### Community 11 - "policyEditor"
Cohesion: 0.12
Nodes (13): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), eventsOrAll(), netWarnings() (+5 more)

### Community 12 - "icon"
Cohesion: 0.06
Nodes (33): backendLink, configView, forwardForm, forwardView, homePage, protectedPage, Event, newLogo() (+25 more)

### Community 13 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.06
Nodes (22): diffModel, editorModel, listenForDaemonEvents(), tickDaemonStats(), formatType(), listenForEvents(), syncViewport(), tickStats() (+14 more)

### Community 14 - ".launch"
Cohesion: 0.27
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 15 - "User"
Cohesion: 0.05
Nodes (69): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir() (+61 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (6): eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (58): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+50 more)

### Community 18 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (12): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, dirEntry (+4 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.15
Nodes (47): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+39 more)

### Community 22 - "Resource"
Cohesion: 0.07
Nodes (12): Resource, GuardRepository, Vault, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots(), TestPartitionEncryptionRootsSplitsDirsAndFiles() (+4 more)

### Community 23 - "setpassword_authz_test.go"
Cohesion: 0.12
Nodes (37): CheckPassword(), Hash(), HashConfigured(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf(), parseHash() (+29 more)

### Community 24 - "MountVouch"
Cohesion: 0.05
Nodes (29): GlobChild, statFunc, ensureBtrfsLayout(), deleteInoKeys(), deleteResKeys(), dropSupersededLocked(), PruneSuperseded(), supersedeMaps() (+21 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (72): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+64 more)

### Community 26 - "expandPlaceholders"
Cohesion: 0.14
Nodes (8): symlinkTargetGlob(), expandPlaceholders(), TestExpandWhitelist(), TestExtraWatchPathsForSkipsMissing(), LibDirGlob, TrustGlob, reservedGlobs(), splitTrustGlob()

### Community 27 - "time.Time"
Cohesion: 0.08
Nodes (25): LockedError, addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), markReady() (+17 more)

### Community 28 - "NewDumpHook"
Cohesion: 0.14
Nodes (15): applyVerbosity(), VerboseLevel, BridgeStdLog(), formatEntry(), NewDumpHook(), sortedKeys(), sortStrings(), TestBridgeStdLog() (+7 more)

### Community 29 - "hashFileShared"
Cohesion: 0.21
Nodes (9): binaryVerifyState, HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled() (+1 more)

### Community 30 - "gui/install.go"
Cohesion: 0.06
Nodes (43): nextConfig(), installGUI(), needsBinary(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), installedBinary, InstallMode, InstallRequest (+35 more)

### Community 31 - "serve.go"
Cohesion: 0.07
Nodes (35): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+27 more)

### Community 32 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 33 - "netModel"
Cohesion: 0.08
Nodes (24): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+16 more)

### Community 34 - "picker"
Cohesion: 0.10
Nodes (18): keyEntry, listOpts, pickEntry, picker, treeSearch, errText(), isFindShortcut(), newKeyEntry() (+10 more)

### Community 35 - "github.com/spf13/cobra.Command"
Cohesion: 0.11
Nodes (24): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+16 more)

### Community 37 - "backend"
Cohesion: 0.14
Nodes (7): checkSecret(), backend, backendOp, backendState, Message, Request, RequireDaemonStopped()

### Community 38 - "controlServer"
Cohesion: 0.08
Nodes (25): readConfigPut(), TestReadGrantRequestConfig(), newControlManager(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest() (+17 more)

### Community 39 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 40 - "os.File"
Cohesion: 0.05
Nodes (57): multicallAdmissible(), ownerPaths(), setConfinedHomes(), vetDecision, VettedInode, Multicall, MulticallKind, Source (+49 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.07
Nodes (15): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+7 more)

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (6): shQuote(), IntegrationSuite, inNosuidNamespace(), IntegrationSuite, IntegrationSuite, absPath()

### Community 44 - "NetGuardEvent"
Cohesion: 0.18
Nodes (6): snapNetGuard, NetGuardEvent, protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 45 - "GuardEvent"
Cohesion: 0.06
Nodes (25): holdHeadless(), HeadlessLine(), runGuardHeadless(), bpfGuardEvent, snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestLocalizePaths() (+17 more)

### Community 46 - "Open"
Cohesion: 0.18
Nodes (19): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+11 more)

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "editView"
Cohesion: 0.11
Nodes (10): codeEntry, editView, Entry, selectable(), fileSize(), joinRel(), newCodeEntry(), newEditView() (+2 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (27): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+19 more)

### Community 50 - "LibraryClosure"
Cohesion: 0.07
Nodes (30): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+22 more)

### Community 51 - "GenerateConf"
Cohesion: 0.04
Nodes (73): logRefreshChange(), TestLogRefreshChange(), appendSectionsAndEdit(), TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), askEncryption(), editConfig(), libraryBlocksFromCandidates() (+65 more)

### Community 52 - "proto.go"
Cohesion: 0.07
Nodes (24): Client, File, User, childEnv(), launch(), Start(), helper(), TestChildEnvIsAllowlisted() (+16 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.10
Nodes (47): systemPatterns(), systemRule(), grantNewBinaries(), watchGroup, Line(), Lines(), TestLinesKeyOnLinkAndIncludePending(), addResource() (+39 more)

### Community 54 - "tabs.go"
Cohesion: 0.03
Nodes (22): appTheme, binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panel (+14 more)

### Community 55 - "TrustGuard"
Cohesion: 0.08
Nodes (11): exeRow, tempGrant, tempMaskOp, tempRow, trustHook, trustRows, setSupersedeTrust(), liveTrust() (+3 more)

### Community 56 - "wizard.go"
Cohesion: 0.07
Nodes (36): releaseBuild(), Run(), deletePostBackups(), listBackups(), restoreBackups(), askDecrypt(), confirmKeyDeletion(), decryptDirectories() (+28 more)

### Community 57 - "Network"
Cohesion: 0.10
Nodes (11): binaryList, newBinaryList(), netGuardPill(), newExeResolver(), guardRow(), clonePolicy(), NetExtra, NetPolicy (+3 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.04
Nodes (45): engine, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates() (+37 more)

### Community 59 - "filevault.go"
Cohesion: 0.05
Nodes (63): deprovisionKind, harnessSuite, migrationTarget, classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath() (+55 more)

### Community 60 - "update/update.go"
Cohesion: 0.05
Nodes (86): installHooks(), staged(), changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), showChangelog(), applyUpdate() (+78 more)

### Community 61 - "newTestTarget"
Cohesion: 0.23
Nodes (16): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+8 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.21
Nodes (15): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+7 more)

### Community 63 - "NetEventType"
Cohesion: 0.09
Nodes (12): snapNetMonitor, NetBpfEvent, NetEvent, NetEventType, newNetEventChecks(), NetEventTypes(), ParseNetEventType(), NetworkMonitorRepository (+4 more)

### Community 65 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 66 - "gui/session.go"
Cohesion: 0.07
Nodes (41): ancestorEnv(), cleanupSessionEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField() (+33 more)

### Community 67 - "runMonitor"
Cohesion: 0.11
Nodes (18): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+10 more)

### Community 68 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 69 - "codeedit.go"
Cohesion: 0.09
Nodes (39): class, classOf(), Edit(), leadingIndent(), New(), onlySpaces(), SetText(), softTabBack() (+31 more)

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "runOfflineEdit"
Cohesion: 0.07
Nodes (35): auditAfterEdit(), auditFindings(), reserveStdout(), runBackend(), dialLiveSession(), LiveModeAvailable(), applyEdit(), checkAndApply() (+27 more)

### Community 72 - "Guard"
Cohesion: 0.06
Nodes (23): deferredBinary, rootHandle, BinaryEntry, openRoot(), TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), checkBtrfsKeys(), SharedPinDegraded(), TestNoteRootMountFollowsMountID() (+15 more)

### Community 73 - "sync.Mutex"
Cohesion: 0.07
Nodes (20): filesPage, Target, specFromForm(), guardPill(), specSummary(), EventTypes(), newActivity(), Files (+12 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.17
Nodes (9): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), btrfsTypes(), TestBtrfsLayoutFromNestedMembers() (+1 more)

### Community 75 - "hashCache"
Cohesion: 0.19
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 76 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 77 - "RefreshCatalog"
Cohesion: 0.16
Nodes (15): sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), refreshCatalogConfig(), RefreshCatalog() (+7 more)

### Community 78 - "prepareDaemonStart"
Cohesion: 0.10
Nodes (27): CheckBPFLSM(), configureDaemonLogging(), loadDaemonConfig(), lockOneRoot(), lockRootRecovering(), newPinGeneration(), prepareDaemonStart(), relockStaleVaults() (+19 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "Monitor"
Cohesion: 0.15
Nodes (4): NewMonitor(), newMountTable(), Monitor, watchedInode

### Community 81 - "updatePage"
Cohesion: 0.14
Nodes (12): latestRelease, updatePage, probeSystem(), boolIndex(), isPreRelease(), Check(), collectDiagnostics(), readSysctl() (+4 more)

### Community 82 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 83 - "dialog_test.go"
Cohesion: 0.36
Nodes (12): pickMode, openPicker(), pickerTree(), shownNames(), TestPickerArrowHighlightsFolderEnterOpensIt(), TestPickerCtrlFFocusesSearch(), TestPickerFolderModeChoosesShownFolder(), TestPickerSearchFiltersFolder() (+4 more)

### Community 84 - ".posAt"
Cohesion: 0.14
Nodes (11): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), setGutter(), stylesOf() (+3 more)

### Community 86 - "find_test.go"
Cohesion: 0.30
Nodes (14): findAll(), lowerRunes(), NewFinder(), runesEqual(), findEditor(), keys(), TestFindAllSmartCase(), TestFinderLeavesSaveToHost() (+6 more)

### Community 87 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 88 - "liveSession"
Cohesion: 0.16
Nodes (13): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), forwardCheck() (+5 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "DaemonEvent"
Cohesion: 0.14
Nodes (24): ServeConfig, catchLifecycleSignals(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI(), runTUI(), admitEvent() (+16 more)

### Community 92 - ".FilterExistingWhitelist"
Cohesion: 0.10
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.09
Nodes (11): foreignPinnedLink(), eventTypeKey(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs() (+3 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "debouncer"
Cohesion: 0.39
Nodes (3): TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), debouncer

### Community 100 - "github.com/stretchr/testify/suite.Suite"
Cohesion: 0.20
Nodes (3): mountinfoUnitTest, unescapeMountinfo(), editorHarness

### Community 101 - "writeEvent"
Cohesion: 0.33
Nodes (8): UIDResolver, NewUIDResolver(), writeEvent(), allowedEvent(), deniedEvent(), TestWriteEventAllowed(), TestWriteEventAllowedSuppressed(), TestWriteEventBlocked()

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "SuperblockDev"
Cohesion: 0.29
Nodes (7): FirstOnBtrfs(), mountinfoDev(), OnBtrfs(), SuperblockDev(), SuperblockDevPath(), TestMountinfoDev(), TestSuperblockDevMatchesStatOffBtrfs()

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "BinaryStat"
Cohesion: 0.12
Nodes (19): BinaryStat, placedWalk, binaryStatOf(), binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf() (+11 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "installPage"
Cohesion: 0.11
Nodes (7): actionGate, binSource, installAction, installPage, uninstallCard, daemonStopNote(), Pulse

### Community 111 - "mountTable"
Cohesion: 0.27
Nodes (7): MountInfo, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), mountTable

### Community 112 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.24
Nodes (11): NewGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle() (+3 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "CandidateDir"
Cohesion: 0.09
Nodes (31): entryWriters(), reserveChain(), TestGlobBuilder_LibBitsPerEntry(), underResource(), selectAndEditConfig(), addManualDirectories(), groupCandidates(), pickDirectories() (+23 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "go_pkg_fyne_io_fyne_v2"
Cohesion: 0.05
Nodes (32): buttonCursor, fakeLink, multiAnswer, pageDef, textAnswer, Hello, handCursors(), handCursorsOn() (+24 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "fakeGuardRepo"
Cohesion: 0.06
Nodes (18): SupersededBinary, TempBinary, TemporaryGrant, TempRule, planTempBlock(), applyTempGrants(), closeTempBinaries(), TemporaryAccess (+10 more)

### Community 143 - "time.Duration"
Cohesion: 0.17
Nodes (12): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+4 more)

### Community 145 - "ParseEventsFlag"
Cohesion: 0.29
Nodes (6): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule()

### Community 146 - "Config"
Cohesion: 0.06
Nodes (41): drain(), newCatalogRefresher(), readEvents(), swapFile(), makeReloadHandler(), noteReloadedSettings(), startCatalogRefresh(), inspectorPaths() (+33 more)

### Community 147 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 148 - "Execute"
Cohesion: 0.29
Nodes (7): Requested(), TestRequested(), TestSubcommandOutranksGUI(), Execute(), initLogger(), subcommand(), main()

### Community 149 - "fscrypt.go"
Cohesion: 0.06
Nodes (47): confirmMasterKeyOverwrite(), runGenKey(), confirmRunPrereq(), resolveFilesystemPrereqs(), cleanOrphanedMetadata(), checkKeyLen(), classifySetupError(), classifySupportError() (+39 more)

### Community 153 - "ConfinedEntry"
Cohesion: 0.26
Nodes (20): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+12 more)

### Community 154 - "EventType"
Cohesion: 0.11
Nodes (10): snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), Cstr() (+2 more)

### Community 155 - "layout.go"
Cohesion: 0.23
Nodes (10): adaptiveGrid, checkGrid, gapVBox, TestAdaptiveGridDropsColumns(), cellMin(), newAdaptiveGrid(), newCheckGrid(), newGapVBox() (+2 more)

### Community 156 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 157 - "journal"
Cohesion: 0.14
Nodes (6): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "backend_test.go"
Cohesion: 0.28
Nodes (18): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+10 more)

### Community 161 - "vaultfs.go"
Cohesion: 0.07
Nodes (35): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+27 more)

### Community 162 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 164 - "pinstate.go"
Cohesion: 0.18
Nodes (19): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID(), TestPinStateRoundTrip() (+11 more)

### Community 165 - "feed.go"
Cohesion: 0.06
Nodes (43): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+35 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (39): installUI, pageID, reel, resizeWatch, runHost, shell, sizedDialog, snapEngines (+31 more)

### Community 167 - "runView"
Cohesion: 0.09
Nodes (9): logGrid, logHook, logLine, runView, stageItem, stageState, logRow(), logStyle() (+1 more)

### Community 168 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 170 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "KernelDev"
Cohesion: 0.09
Nodes (26): fdStat, physicalParent(), statFD(), isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed() (+18 more)

### Community 179 - "install/install.go"
Cohesion: 0.04
Nodes (96): applyLiveRefresh(), buildBinaryIfNeeded(), certifyInstall(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), encryptDirectories(), ensureInstalledBinary(), ensureMasterKey() (+88 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 189 - "NetworkMonitor"
Cohesion: 0.10
Nodes (11): FormatAddr(), Ntohs(), WithSuspectsFrom(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), RealNetEngines() (+3 more)

### Community 190 - "Development"
Cohesion: 0.29
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

### Community 196 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 204 - "runNetworkGuard"
Cohesion: 0.09
Nodes (24): CheckEBPF(), ParseNetEventsFlag(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard() (+16 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.15
Nodes (13): Actions, Desktop GUI, File and folder pickers, Files, Home, Install, Network, Protected (+5 more)

### Community 209 - "Ledger"
Cohesion: 0.12
Nodes (11): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), fyneLogFilter, applySchema() (+3 more)

### Community 210 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 216 - "binaryVetter"
Cohesion: 0.05
Nodes (63): newBinaryVetter(), openBinaryVetter(), backingDeviceUnion(), buildConcurrency(), buildGuards(), buildOneGuard(), drainEphemeral(), eventFilterOptions() (+55 more)

### Community 222 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

## Knowledge Gaps
- **31 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+26 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 566 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **32 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _31 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.12224658106274142 - nodes in this community are weakly interconnected._
- **Why does `Editor keys` connect `edit-protected` to `liveSession`, `wizard.go`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.024287856071964018 - nodes in this community are weakly interconnected._
- **Why does `Demo GIFs` connect `Development` to `shell`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08957708957708958 - nodes in this community are weakly interconnected._